# Vidura Sense Design Specification

## Table of Contents

[toc]

## Rivision History

| Version | Authors | Description   | Date       |
| ------- | ------- | ------------- | ---------- |
| 1.0     | Bikash  | Initial Draft | 2023-07-17 |

## 1.0 Introduction

The purpose of this document is to specify the requirements for an intelligent web crawler bot.
The crawler will crawl the World Wide Web to gather info.

## 2.0 Requirements

### 2.1 Requester

This module will have the apis to call urls and fetch content.

- Call [Frontier](#22-forntier) for urls to get content form World Wide Web.
- Call multiple urls based upon their preferance provided by frontier concurrently.
- Pass web contents to the [Data Processor](#23-data-processor)

### 2.2 Forntier

This contains the apis to feed urls for every trigger.

- Fetch urls from database whose content is 3(get it from configuration "update_interval") days old or more (ie current_date - last_updated >=3 days) and return as array.

### 2.3 Data Processor

This contains the apis to process data.

- Get content form [Requestor](#21-requester).
- Process it to schema refer to [WebCache](./Schemas.md#10-webcache)
- Filter content using Intelligence provider APIs.
- After processing send it to Storage Provider

### 2.4 Storage Provider

This module contains apis to store data in databases or other storages.

- Store content in database received from [Data Processor](#23-data-processor)

### 2.5 Triggerer

This module triggers [Requestor](#21-requester) on regular intervals.

### 2.6 REST Apis

This module contains the restapis for triggered jobs,metrics and web content summary.

## 3.0 System Disgram

![image](diagrams/SystemDiagram.svg)