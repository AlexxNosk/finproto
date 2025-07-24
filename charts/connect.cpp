#include "connect.h"


extern "C" {
    void RunGoFunction(const char* token);  // defined in Go
}

void loadEnvFile(const std::string& path) {
    std::ifstream file(path);
    if (!file.is_open()) {
        std::cerr << "Could not open .env file\n";
        return;
    }

    std::string line;
    while (std::getline(file, line)) {
        auto delimiterPos = line.find('=');
        if (delimiterPos == std::string::npos) continue;
        std::string key = line.substr(0, delimiterPos);
        std::string value = line.substr(delimiterPos + 1);

        // Remove quotes if any
        if (value.front() == '"' && value.back() == '"') {
            value = value.substr(1, value.size() - 2);
        }

        setenv(key.c_str(), value.c_str(), 1);  // overwrite = 1
        std::cout << "KLUCH: " << key << " = " << value << std::endl;
    }
}
