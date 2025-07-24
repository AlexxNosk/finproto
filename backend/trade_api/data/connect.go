package data

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alexxnosk/finproto/backend/trade_api/v1/accounts/accounts_service"
	"github.com/alexxnosk/finproto/backend/trade_api/v1/assets/assets_service"
	"github.com/alexxnosk/finproto/backend/trade_api/v1/auth/auth_service"
	"github.com/alexxnosk/finproto/backend/trade_api/v1/marketdata/marketdata_service"
	"github.com/alexxnosk/finproto/backend/trade_api/v1/orders/orders_service"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

const (
	endPoint = "api.finam.ru:443" //"ftrr01.finam.ru:443"
	connPgStr = "postgres://root:root@localhost:5434/finProto_db"
)

var logLevel = &slog.LevelVar{} // INFO
var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: logLevel,
})).With(slog.String("package", "go-finam-grpc"))

func SetLogger(logger *slog.Logger) {
	log = logger
}

// SetLogDebug установим уровень логгирования Debug
func SetLogDebug(debug bool) {
	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
}

type Asset struct {
	Ticker string `json:"ticker,omitempty"` // Тикер инструмента
	Symbol string `json:"symbol,omitempty"` // Символ инструмента ticker@mic
	Name   string `json:"name,omitempty"`   // Наименование инструмента
	Mic    string `json:"mic,omitempty"`    // Идентификатор биржи
	Type   string `json:"type,omitempty"`   // Тип инструмента
	Id     string `json:"id,omitempty"`     // Идентификатор инструмента
}

type Client struct {
	token             string    // Основой токен пользователя
	accessToken       string    // JWT токен для дальнейшей авторизации
	ttlJWT            time.Time // Время завершения действия JWT токена
	conn              *grpc.ClientConn
	connPG			  *pgx.Conn
	AuthService       auth_service.AuthServiceClient
	AccountsService   accounts_service.AccountsServiceClient
	AssetsService     assets_service.AssetsServiceClient
	MarketDataService marketdata_service.MarketDataServiceClient
	OrdersService     orders_service.OrdersServiceClient
	Securities        map[string]Asset //  Список инструментов с которыми работаем (или весь список? )
}

func NewClient(ctx context.Context, token string) (*Client, error) {
	// TODO выделить в отдельный метод connect()
	log.Debug("NewClient start connect")
	conn, err := grpc.NewClient(endPoint,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                15 * time.Minute, // отправлять ping каждые 15 минут
			Timeout:             30 * time.Second, // ждать ответа не дольше 10 сек
			PermitWithoutStream: true,             // пинговать даже без активных RPC
		}),
	)
	if err != nil {
		return nil, err
	}
	connPG, err:= ConnPG(ctx, connPgStr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		token:             token,
		conn:              conn,
		connPG: 		   connPG,
		AuthService:       auth_service.NewAuthServiceClient(conn),
		AccountsService:   accounts_service.NewAccountsServiceClient(conn),
		AssetsService:     assets_service.NewAssetsServiceClient(conn),
		MarketDataService: marketdata_service.NewMarketDataServiceClient(conn),
		OrdersService:     orders_service.NewOrdersServiceClient(conn),
		Securities:        make(map[string]Asset),
	}
	log.Debug("NewClient есть connect")
	err = client.UpdateJWT(ctx) // сразу получим и запишем accessToken для работы
	if err != nil {
		return nil, err
	}
	// в отдельном потоке периодически обновляем accessToken
	go client.runJwtRefresher(ctx)
	return client, nil
}

func (c *Client) Close(ctx context.Context) error {
	err := c.connPG.Close(ctx)
	if err !=nil{
		return err
	}
	err = c.conn.Close()
	if err !=nil{
		return err
	}
	return nil

}

func ConnPG (ctx context.Context, connPgStr string) (*pgx.Conn, error){
	connPG, err := pgx.Connect(ctx, connPgStr)
	if err != nil {
		slog.Error("pgx.Connect", "err", err.Error())
		return nil, err
	}
    return connPG, err
}







const authKey = "Authorization"             //
const jwtTokenTtl = 12 * time.Minute        // Время жизни токена JWT в минутах (15 минут)
const jwtRefreshInterval = 10 * time.Minute // Интервал обновления токена (в минутах)

// GetJWT Получение JWT токена из API токена
//
// идет  вызов AuthService.Auth
func (c *Client) GetJWT(ctx context.Context) (string, error) {
	if c.token == "" {
		c.accessToken = ""
		return c.accessToken, nil
	}
	req := &auth_service.AuthRequest{Secret: c.token}
	log.Debug("GetJWT start AuthService.Auth")
	t := time.Now()
	res, err := c.AuthService.Auth(ctx, req)
	if err != nil {
		return c.accessToken, err
	}
	log.Debug("GetJWT end AuthService.Auth", "duration", time.Since(t))
	return res.Token, nil
}

// UpdateJWT
// если jwt токен пустой или вышло его время
// получим JWT
// и запишем его в параметры клиента. Проставим время получения
func (c *Client) UpdateJWT(ctx context.Context) error {
	if c.token == "" {
		c.accessToken = ""
		return fmt.Errorf("UpdateJWT: token пустой")
	}

	if c.accessToken == "" || c.ttlJWT.Before(time.Now()) {
		token, err := c.GetJWT(ctx)
		if err != nil {
			log.Error("UpdateJWT: failed to refresh JW", "err", err.Error())
			return err
		}
		c.ttlJWT = time.Now().Add(jwtTokenTtl)
		c.accessToken = token
	}
	return nil
}

// runJwtRefresher в отдельном потоке периодически обновляем токен
func (c *Client) runJwtRefresher(ctx context.Context) {
	log.Debug("run JwtRefresher")
	if c.token == "" {
		c.accessToken = ""
		log.Warn("JWT refresher token пустой. Выход из метода")
		return
	}
	ticker := time.NewTicker(jwtRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debug("start JwtRefresh")
			token, err := c.GetJWT(ctx)
			if err != nil {
				log.Error("JwtRefresher: failed to refresh JW", "err", err.Error())
			}
			// запишем время окончания токена
			c.ttlJWT = time.Now().Add(jwtTokenTtl)
			c.accessToken = token
		case <-ctx.Done():
			log.Debug("JWT refresher stopped")
			return
		}
	}
}

// GetTokenDetails Получение информации о токене сессии
//
// идет вызов AuthService.TokenDetails
func (c *Client) GetTokenDetails(ctx context.Context) (*auth_service.TokenDetailsResponse, error) {
	return c.AuthService.TokenDetails(ctx, &auth_service.TokenDetailsRequest{Token: c.accessToken})
}

// WithAuthToken Создаем новый контекст с заголовком Authorization
// пишем в него jwt токен
func (c *Client) WithAuthToken(ctx context.Context) (context.Context, error) {
	// проверим наличие токена
	err := c.UpdateJWT(ctx)
	if err != nil {
		return ctx, err
	}

	//_ = authKey
	//ctx = context.WithValue(ctx, authKey, c.accessToken)
	//return ctx, nil
	// добавим заголовок
	md := metadata.New(map[string]string{
		authKey: c.accessToken,
	})
	// и добавляем его в ctx
	return metadata.NewOutgoingContext(ctx, md), nil
}
