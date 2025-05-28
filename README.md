# 🐾 ArmyK9 - AI-Developed Cybersecurity Tools

Welcome to **ArmyK9**, a collection of AI-developed tools designed to automate various aspects of cybersecurity monitoring. Whether you're a 🛡️ penetration tester, 🧪 security researcher, or an 💼 IT professional, these tools can assist in tracking vulnerabilities and gaining a clearer understanding of your organization's security posture.

## 📦 Overview

This repository contains scripts and utilities that help automate the process of collecting exploit information, making it easier for organizations to stay informed about relevant vulnerabilities. The primary focus is on:

- 🐞 **Vulnerability listing**
- 🧭 **Assisting in determining security posture**
- 📊 **Data collection for vulnerability assessments**

## 📁 Key Files

- **`whatmatters.go`** 🧠: Automates the collection of Remote Code Execution (RCE) exploits from the past 12 months, using RSS feeds from popular exploit code publishers like ExploitDB and PacketStormSecurity. It downloads relevant exploit code into a folder, making it easier to track emerging threats.

- **`exploits_list.go`** 🖥️: Generates a 1-page HTML listing of the retrieved exploits, organized into a clean table format. The report includes clickable links to detailed information for each exploit, allowing quick review of vulnerabilities published over a given period.

## 🎯 Purpose

These tools were designed with a dual-purpose mindset:

- 🕵️ On the **offensive** side, they help penetration testers automate the collection of relevant exploit code.  
- 🛡️ On the **defensive** side, they allow organizations to keep track of recently published RCE exploits and assess whether their systems are vulnerable.

The tools are especially useful for smaller organizations with 🪙 limited budgets for high-end cybersecurity solutions. By automating exploit collection, they help maintain an up-to-date view of potential vulnerabilities without expensive infrastructure.

## 🛠️ Installation Guide for WhatMatters

**Important Notes**:
- ⏱️ This tool downloads only RCE-related exploit codes from the last 12 months.
- 📂 Files are renamed based on category, title, and programming language.
- ⚠️ Some files may contain single-line payloads and may not compile OOTB.
- 🧠 Use ChatGPT to help fix any non-compiling code before using it.

### 🧬 Steps to Install and Compile:

1. **Clone the repository**:
   ```bash
   git clone https://github.com/armyk9/whatmatters.git
   ```

2. **Initialize the Go Module**:
   ```bash
   go mod init whatmatters
   ```

3. **Install GoQuery**:
   ```bash
   go get github.com/PuerkitoBio/goquery
   ```

4. **Install GoFeed**:
   ```bash
   go get github.com/mmcdole/gofeed
   ```

5. **Compile the Application**:

   #### 🪟 For Windows:
   ```bash
   go build -o whatmatters.exe whatmatters.go
   go build -o exploits_list.exe exploits_list.go
   ```

   #### 🐧 For Linux:
   ```bash
   go build -o whatmatters whatmatters.go
   go build -o exploits_list exploits_list.go
   ```

## ❓ Why WhatMatters?

A friend of mine, who serves as a Cybersecurity Manager at a local bank 🏦, recently shared his concerns about the lack of timely visibility into newly published exploits. To address this, I created a tool that scrapes exploit databases and delivers real-time insights into potential threats 🚨 — powered by ChatGPT.

## 📜 License

This project is licensed under the MIT License. See the LICENSE file for details.
