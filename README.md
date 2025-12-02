# gRPC CRUD Application

풀스택 gRPC 애플리케이션입니다. Go 서버와 React 클라이언트를 사용합니다.

## 기술 스택

### Backend
- Go
- Connect (gRPC)
- Buf

### Frontend
- Vite
- React
- TanStack Router
- TailwindCSS
- Connect-Web

## 프로젝트 구조

```
.
├── proto/                      # Protocol Buffer 정의
│   └── item/v1/
│       └── item.proto
├── server/                     # Go 서버
│   ├── proto-generated/        # 생성된 Go 코드
│   ├── main.go
│   └── go.mod
├── client/                     # React 클라이언트
│   ├── src/
│   │   └── proto-generated/    # 생성된 TypeScript 코드
│   └── package.json
├── scripts/                    # 유틸리티 스크립트
│   ├── generate.sh             # 코드 생성
│   └── dev.sh                  # 개발 서버 실행
├── buf.yaml
└── buf.gen.yaml
```

## 시작하기

### 필수 요구사항

- Go 1.23+
- Node.js 18+
- Buf CLI

Buf 설치:
```bash
# macOS
brew install bufbuild/buf/buf

# 기타
npm install -g @bufbuild/buf
```

### 설정

1. 코드 생성:
```bash
./scripts/generate.sh
```

2. Go 의존성 설치:
```bash
cd server
go mod download
cd ..
```

3. Node 의존성 설치:
```bash
cd client
npm install
cd ..
```

### 개발 서버 실행

두 서버를 동시에 실행:
```bash
./scripts/dev.sh
```

또는 개별 실행:

**Go 서버:**
```bash
cd server
go run main.go
# 서버: http://localhost:8080
```

**React 클라이언트:**
```bash
cd client
npm run dev
# 클라이언트: http://localhost:5173
```

## API 엔드포인트

- `CreateItem` - 아이템 생성
- `GetItem` - 아이템 조회
- `ListItems` - 아이템 목록 조회
- `UpdateItem` - 아이템 수정
- `DeleteItem` - 아이템 삭제

## 코드 생성

proto 파일 수정 후:
```bash
./scripts/generate.sh
```

이 스크립트는 다음을 생성합니다:
- Go 코드: `server/proto-generated/item/v1/`
- TypeScript 코드: `client/src/proto-generated/item/v1/`
