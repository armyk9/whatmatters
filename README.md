# ArmyK9 - AI-Developed Cybersecurity Tools

Welcome to **ArmyK9**, a collection of AI-developed tools designed to automate various aspects of cybersecurity monitoring. Whether you're a penetration tester, security researcher, or an IT professional, these tools can assist in tracking vulnerabilities and gaining a clearer understanding of your organization's security posture.

## Overview

This repository contains scripts and utilities that help automate the process of collecting exploit information, making it easier for organizations to stay informed about relevant vulnerabilities. The primary focus is on:

- **Vulnerability listing**
- **Assisting in determining security posture**
- **Data collection for vulnerability assessments**

## Key Files

- **`whatmatters.go`**: Automates the collection of Remote Code Execution (RCE) exploits from the past 12 months, using RSS feeds from popular exploit code publishers like ExploitDB and PacketStormSecurity. It downloads relevant exploit code into a folder, making it easier to track emerging threats.

- **`vulns_report.go`**: Generates a 1-page HTML listing of the retrieved exploits, organized into a clean table format. The report includes clickable links to detailed information for each exploit, allowing quick review of vulnerabilities published over a given period.

## Purpose

These tools were designed with a dual-purpose mindset. On the offensive side, they help penetration testers automate the collection of relevant exploit code. On the defensive side, they allow organizations to keep track of recently published RCE exploits and assess whether their operating systems, software, or platforms are vulnerable. This helps organizations focus their efforts on patching or upgrading systems accordingly, leading to a better overall security posture.

The tools are particularly useful for smaller organizations with limited budgets for high-end cybersecurity solutions. By automating exploit collection, they help maintain an up-to-date view of potential vulnerabilities without needing expensive infrastructure.

## Installation Guide for WhatMatters

**Important Notes**:
- This tool downloads only RCE-related exploit codes from the last 12 months.
- All downloaded files are renamed based on their category, title, and programming language.
- Some text files may not be directly compilable but contain single-line payloads that can still be used.
- Not all downloaded exploit code will compile or work out-of-the-box (OOTB). Use ChatGPT to help correct and fix any issues before compiling or executing the code.

### Steps to Install and Compile:

1. **Clone the repository**:
   
   ```bash
   git clone https://github.com/derang3d/ArmyK9

2. **Initialize the Go Module**:  
   Run the following command in your CLI, within the directory where `whatmatters.go` is located:

   ```bash
   go mod init whatmatters

3. **Install GoQuery**:   
   Install GoQuery, which is required for parsing web content, by running the following command:

   ```bash
   go get github.com/PuerkitoBio/goquery

4. **Install GoFeed**:
   Install GoFeed for parsing RSS feeds by running the command:

   ```bash
   go get github.com/mmcdole/gofeed

5. **Compile the Application**:
   Once the dependencies are installed, compile the application by running this command in the build directory:

   ```bash
   go build -o whatmatters.exe whatmatters.go

# Why ArmyK9?
The name ArmyK9 reflects my background in both the military and cybersecurity fields. Just as K9 units are deployed for a variety of tasks, this project aims to provide practical tools for cybersecurity tasks in today's challenging landscape. These AI-powered scripts are designed to enhance security workflows and provide a better understanding of an organization's security posture.

# License
This project is licensed under the MIT License. See the LICENSE file for details.
