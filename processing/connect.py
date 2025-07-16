import psycopg2
import pandas as pd
import plotly.graph_objects as go

# Connect to PostgreSQL
conn = psycopg2.connect(
    host="localhost",
    port=5434,
    database="finProto_db",
    user="root",
    password="root"
)

	# symbol := "VTBR@MISX"
	# tfStr := "D"
	# startStr := "01-01-2000"
	# endStr := "nil"
	# action := ""
	# type DataStruct struct{
	# 		X float64
	# 		Y float64
	# 		Z float64
	# } 





# Load bars table into DataFrame
df = pd.read_sql("SELECT * FROM bars_vtbr_misx_d ORDER BY timestamp", conn)
conn.close()

# Make sure timestamp is datetime
df['timestamp'] = pd.to_datetime(df['timestamp'])

# Plot OHLC candlestick chart
fig = go.Figure(data=[
    go.Candlestick(
        x=df['timestamp'],
        open=df['open'],
        high=df['high'],
        low=df['low'],
        close=df['close'],
        name="OHLC"
    )
])

fig.update_layout(title="OHLC Chart", xaxis_title="Time", yaxis_title="Price")
fig.show()
