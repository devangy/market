# tradeFlow-v1


**tradeFlow-v1** is an Telegram Bot that detects whales or smart money(insiders, traders with good win rate, Arbitrage bots) on Prediction markets like Polymarket and Kalsi. It does so by calculating metrics on the past 50-100 trades of a user  that has traded or bought a event yes or no share for. The bot calculates winrate, max drawdown, wins, losses, pnl, net pnl etc and flags arbitrage bots that buy and sell fast profiting from volatility. The bot polls 4 apis to get data data from and check for fresh events popping up on the market deuplicates using a hashmap and writes to a small binary file which contains hashes of seen events for deuplicaton. The bot is running on aws ec2 instance created with terraform (IaC) with Docker compose.


What is a Prediction Market ? 

> On prediction markets we buy shares for an event. There are 2 types of shares we can buy and in any quantity
> Yes or No
> The higher the chances of something happening the more the share price
> example:
> Event Name: Will India win the world Cup in 2026?
> I buy 10 YES shares for 0.50$ cents because I think the chances are 80 percent of it happening.
> The market resolves to 1.00 $ always
> so I buy 10 shares at 0.50  = 5 $
> market resolves to yes I win
> I get back 10 $ direct 2x 

---

## Live Demo
You can view the real-time status and performance dashboard of the bot here:
Telegram Handle ->   @datalog02bot

---

## Features
* **Multi-Exchange Support:** Simultaneous integration with Kalshi and Polymarket.
* **Signal Confidence Scoring:** Uses weighted indicators (RSI, EMA, MACD) to filter noise.
* **Dockerized Environment:** Deploy anywhere with a single command.
* **CI/CD Integration:** Automated testing and deployment via GitHub Actions.

---

## Tech Stack
* **Runtime:** Go 1.24.7
* **Infrastructure:** Docker and docker compose
* **CI/CD:** GitHub Actions
* **APIs:** Kalshi API, Polymarket API

---

## Installation & Setup

### Prerequisites
* Docker and Docker Compose installed.
* API Keys for the respective exchanges.

### 1. Clone the Repository
```bash
git clone [https://github.com/devangy/tradeFlow-v1.git](https://github.com/devangy/tradeFlow-v1.git)
cd tradeFlow-v1


[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
