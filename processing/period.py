import numpy as np
import pandas as pd
from sqlalchemy import create_engine, text
import plotly.graph_objects as go
from plotly.subplots import make_subplots
import matplotlib.pyplot as plt
from statsmodels.graphics.tsaplots import plot_acf
# --- 1. Database connection ---
db_config = {
    'user': 'root',
    'password': 'root',
    'host': 'localhost',
    'port': '5434',
    'database': 'finProto_db'
}

engine = create_engine(
    f"postgresql://{db_config['user']}:{db_config['password']}@{db_config['host']}:{db_config['port']}/{db_config['database']}"
)

# --- 2. Load close prices for SMA calculation ---
query_close = "SELECT timestamp, close FROM bars_sber_misx_d ORDER BY timestamp;"
df_close = pd.read_sql(query_close, engine)
df_close['timestamp'] = pd.to_datetime(df_close['timestamp']).dt.floor('s')
df_close.set_index('timestamp', inplace=True)

# --- 3. Calculate SMAs ---
df_close['sma_20'] = df_close['close'].rolling(window=20).mean()
df_close['sma_200'] = df_close['close'].rolling(window=200).mean()
df_close['sma_250'] = df_close['close'].rolling(window=250).mean()

# --- 4. Save to bars_sber_misx_d_data table ---
df_sma = df_close[['sma_20', 'sma_200', 'sma_250']].dropna().copy()
df_sma.reset_index(inplace=True)

with engine.begin() as conn:
    for _, row in df_sma.iterrows():
        conn.execute(
            text("""
                UPDATE bars_sber_misx_d_data
                SET sma_20 = :sma_20,
                    sma_200 = :sma_200,
                    sma_250 = :sma_250
                WHERE timestamp = :ts
            """),
            {
                'sma_20': float(row['sma_20']),
                'sma_200': float(row['sma_200']),
                'sma_250': float(row['sma_250']),
                'ts': row['timestamp']
            }
        )

# --- 5. Load full OHLC + SMA + Volume data for plotting ---
query_full = """
    SELECT d.timestamp, d.open, d.high, d.low, d.close, d.volume,
           da.sma_20, da.sma_200, da.sma_250
    FROM bars_sber_misx_d d
    LEFT JOIN bars_sber_misx_d_data da ON d.timestamp = da.timestamp
    ORDER BY d.timestamp;
"""
df = pd.read_sql(query_full, engine)
df['timestamp'] = pd.to_datetime(df['timestamp'])
df.set_index('timestamp', inplace=True)
df['close'].plot(title="Close Price Over Time", figsize=(12, 4))
# plt.show()

plot_acf(df['close'].dropna(), lags=100)
plt.title("Autocorrelation")
plt.show()
close = df['close'].dropna()
fft = np.fft.fft(close - close.mean())
frequencies = np.fft.fftfreq(len(fft))

plt.figure(figsize=(10, 4))
plt.plot(frequencies[1:len(frequencies)//2], np.abs(fft)[1:len(frequencies)//2])
plt.title('Frequency Spectrum')
plt.xlabel('Frequency')
plt.ylabel('Amplitude')
plt.show()


# # --- 6. Plot: OHLC + Volume + SMAs ---
# fig = make_subplots(rows=2, cols=1, shared_xaxes=True, 
#                     vertical_spacing=0.05, row_heights=[0.7, 0.3])

# # Candlestick
# fig.add_trace(go.Candlestick(
#     x=df.index,
#     open=df['open'], high=df['high'],
#     low=df['low'], close=df['close'],
#     name='OHLC'), row=1, col=1)

# # SMAs
# fig.add_trace(go.Scatter(
#     x=df.index, y=df['sma_20'],
#     mode='lines', name='SMA 20', line=dict(color='blue')), row=1, col=1)

# fig.add_trace(go.Scatter(
#     x=df.index, y=df['sma_200'],
#     mode='lines', name='SMA 200', line=dict(color='green')), row=1, col=1)

# fig.add_trace(go.Scatter(
#     x=df.index, y=df['sma_250'],
#     mode='lines', name='SMA 250', line=dict(color='red')), row=1, col=1)

# # Volume
# fig.add_trace(go.Bar(
#     x=df.index, y=df['volume'],
#     name='Volume', marker_color='black'), row=2, col=1)

# fig.update_layout(title='OHLC Chart with SMAs and Volume',
#                   xaxis_rangeslider_visible=False)

# fig.show()
