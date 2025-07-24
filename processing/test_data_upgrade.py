# test_call.py
import ctypes
import json

# Load the shared library
lib = ctypes.cdll.LoadLibrary("./cgo/libinstrument.so")

# Define argument/return types
lib.InstrumentCreateC.argtypes = [ctypes.c_char_p]
lib.InstrumentCreateC.restype = ctypes.c_char_p

# Example JSON payload
# json_data = {
#     "request": {
#         "symbol": "VTBR@MISX",
#         "timeframe": "D",
#         "start_date": "01-01-2000",
#         "end_date": "nil",
#         "operation": "deete"
#     },
#     "tables": {
#         "data": {
#             "sma_20": "float64",
#             "sma_200": "float64",
#             "sma_250": "float64",
#             "ind_1": "float64",
#             "ind_2": "float64",
#             "ind_3": "float64"
#         }
#     }
# }

json_data = {
    "request": {
        "symbol": "VTBR@MISX",
        "timeframe": "D",
        "start_date": "01-01-2000",
        "end_date": "01-01-2000",
        "operation": "create"
    }
}



# Serialize to JSON string and encode as bytes
json_bytes = json.dumps(json_data).encode("utf-8")

# Call the function
result = lib.InstrumentCreateC(json_bytes)

# Print the response
print("Go says:", ctypes.string_at(result).decode("utf-8"))
