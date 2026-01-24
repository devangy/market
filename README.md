# Whale Tracker Telegram Bot

## Overview
This Telegram bot is designed to track whale activities in prediction markets. It provides real-time notifications and analysis of significant trades made by large market participants (whales), aiding users in making informed trading decisions.

## Features
- **Real-time Whale Alerts**: Notifies users of large trades enabling timely responses.
- **Market Analysis**: Provides insights and trends based on whale activity.
- **User Customization**: Users can set their own thresholds for alerts based on trade size and other parameters.
- **Interactive Commands**: Users can interact with the bot through intuitive commands to get updates and insights.

## Getting Started
1. **Requirements**:
   - A Telegram account.
   - Access to the prediction market APIs.
   - Python 3.8 or higher.

2. **Installation**:
   Clone the repository:
   ```bash
   git clone https://github.com/devangy/tradeFlow-v1.git
   cd tradeFlow-v1
   ```
   Install the necessary packages:
   ```bash
   pip install -r requirements.txt
   ```

3. **Configuration**:
   Create a `.env` file in the root directory and add your Telegram Bot token and other necessary API keys.

4. **Running the Bot**:
   Execute the following command:
   ```bash
   python bot.py
   ```

## Usage
- **Start the Bot**: Search for the bot username on Telegram and initiate a chat.
- **Commands**:
  - `/start`: To start interacting with the bot.
  - `/help`: To list available commands and get help.
  - `/set_alert <amount>`: To set a whale alert for trades above a specified amount.

## Contributing
Contributions are welcome! Please submit a pull request to the main branch with your proposed changes.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.