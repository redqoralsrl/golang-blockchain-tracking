# 블록체인 트랜잭션 추적기

## 프로젝트 개요

이 프로젝트는 EVM 블록체인의 트랜잭션을 실시간으로 추적하고 분석하는 시스템입니다. 트랜잭션의 발생, 전파, 확인 과정을 모니터링하고, 관련 데이터를 저장하여 분석할 수 있는 기능을 제공합니다.

- 추가 기능 개선 : init.sql의 balance 잔액이 0 이하일 시 조건은 일시적으로 해제하였습니다.

## 기술 스택

- **Backend**: Go 1.24.2
- **Database**: PostgreSQL
- **Blockchain**: Ethereum (go-ethereum)
- **Container**: Docker & Docker Compose
- **Logging**: zerolog

## 시스템 아키텍처

### 1. 핵심 컴포넌트

- **트랜잭션 모니터링 서비스**: EVM 네트워크의 트랜잭션을 실시간으로 감시
- **데이터 저장소**: PostgreSQL을 사용한 트랜잭션 데이터 저장
- **분석 엔진**: 트랜잭션 패턴 분석 및 통계 처리

### 2. 디렉토리 구조

```
.
├── cmd/            # 메인 애플리케이션 진입점
├── config/         # 설정 파일
├── internal/       # 내부 패키지
├── init.sql        # 데이터베이스 스크립트
├── docker-compose.yaml  # 컨테이너 구성
└── Makefile        # 빌드 및 실행 스크립트
```

### 3. 주요 기능

- 실시간 트랜잭션 모니터링
- 트랜잭션 데이터 저장 및 관리
- 트랜잭션 패턴 분석
- 성능 모니터링 및 로깅

## 시스템 요구사항

- Docker & Docker Compose
- Go 1.24.2 이상
- 최소 2GB RAM
- 2 CPU 코어

## 설치 및 실행

1. 환경 변수 설정

```bash
cp .env.template .env
# .env 파일을 편집하여 필요한 설정 입력
```

2. 데이터베이스 실행

```bash
make local-run
```

3. 애플리케이션 실행

```bash
make run
```

## 성능 최적화

- 트랜잭션 처리 최적화
- 리소스 사용량 제한 설정
- 병렬처리 분석 시스템
