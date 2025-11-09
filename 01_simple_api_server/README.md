# 학습 로그 #01: 간단한 HTTP 서버 만들기

## 학습 목표

Go의 표준 라이브러리 `net/http`를 사용하여 간단한 HTTP 서버를 구현하고, 경로 라우팅 동작을 테스트했습니다.

## 학습 내용

### ServeMux 이해하기

`http.HandleFunc`는 내부적으로 기본 `ServeMux`를 사용합니다. `ServeMux`는 HTTP 요청 멀티플렉서로, 들어오는 요청의 URL을 등록된 패턴 목록과 비교하여 가장 일치하는 패턴의 핸들러를 호출합니다.

### 패턴 매칭 예제 ([참고](https://pkg.go.dev/net/http#ServeMux))

패턴은 요청의 메서드, 호스트, 경로를 매칭할 수 있습니다:

- **`"/index.html"`**: 어떤 호스트와 메서드에서든 경로 `"/index.html"`을 매칭합니다
- **`"GET /static/"`**: 경로가 `"/static/"`으로 시작하는 GET 요청을 매칭합니다
- **`"example.com/"`**: 호스트가 `"example.com"`인 모든 요청을 매칭합니다
- **`"example.com/{$}"`**: 호스트가 `"example.com"`이고 경로가 정확히 `"/"`인 요청을 매칭합니다

### `http.HandleFunc`의 경로 매칭 규칙

- **`"/"`**: 모든 경로를 매칭합니다 (가장 낮은 우선순위)
- **`"/{$}"`**: 정확히 `"/"` 경로만 매칭합니다
- **`"/endpoint1"`**: 정확히 `"/endpoint1"` 경로만 매칭합니다
- **`"/endpoint2/"`**: `"/endpoint2/"`로 시작하는 모든 경로를 매칭합니다

### 슬래시(`/`)로 끝나는 경로의 특별한 동작

- `"/endpoint2"` (슬래시 없음)로 요청하면 → `301 Moved Permanently` 응답 후 `"/endpoint2/"`로 리다이렉트됩니다
- `"/endpoint2/"` (슬래시 있음)로 요청하면 → 정상적으로 핸들러가 실행됩니다

### 테스트 결과

| 요청 경로 | 응답 핸들러 | 설명 |
|---------|-----------|------|
| `/anypath` | HandleFunc #1 | `"/"`가 모든 경로를 매칭 |
| `/` | HandleFunc #2 | `"/{$}"`가 정확히 `"/"`만 매칭 |
| `/endpoint1` | HandleFunc #3 | 정확히 `"/endpoint1"`만 매칭 |
| `/endpoint1/` | HandleFunc #1 | `"/"`가 매칭 (우선순위) |
| `/endpoint2` | 301 리다이렉트 | `"/endpoint2/"`로 자동 리다이렉트 |
| `/endpoint2/` | HandleFunc #4 | `"/endpoint2/"`로 시작하는 경로 매칭 |
| `/endpoint2/anypath` | HandleFunc #4 | `"/endpoint2/"`로 시작하는 경로 매칭 |

### 구현 코드

```go
// 익명 함수를 변수에 할당하여 핸들러 정의
h1 := func(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a HandleFunc #1!\n")
}

// 경로와 핸들러 등록
http.HandleFunc("/", h1)

// 서버 시작
log.Fatal(http.ListenAndServe(":8080", nil))
```

### 코드 개선: `/{$}` 패턴 사용하기

예전에는 정확한 경로 매칭을 위해 핸들러 내부에서 수동으로 경로를 체크했습니다:

```go
// 예전 방식: 요청 정보를 수동으로 체크
h1 := func(w http.ResponseWriter, r *http.Request) {
    // 요청 경로가 정확히 "/" 가 아니면 404 Not Found 응답
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    // 정확히 "/" 일 때만 응답
    io.WriteString(w, "Hello from a HandleFunc #1!\n")
}

http.HandleFunc("/", h1)  // "/"는 모든 경로를 매칭하므로 수동 체크 필요
```

하지만 `/{$}` 패턴을 사용하면 코드가 훨씬 간결해집니다:

```go
// 개선된 방식: 패턴으로 정확한 매칭
h2 := func(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a HandleFunc #2!\n")
}

http.HandleFunc("/{$}", h2)  // 정확히 "/"만 매칭되므로 핸들러 내부 체크 불필요
```

**장점:**
- 핸들러 코드가 간결해집니다
- 경로 매칭 로직이 패턴에 명시적으로 표현됩니다
- 불필요한 조건문과 에러 처리가 제거됩니다

## 학습 정리

1. Go의 HTTP 서버는 표준 라이브러리만으로 간단하게 구현 가능합니다
2. 경로 매칭은 등록 순서와 슬래시 유무에 따라 동작이 달라집니다
3. 슬래시로 끝나는 경로는 해당 경로로 시작하는 모든 하위 경로를 처리합니다
4. `/{$}` 패턴을 사용하면 핸들러 내부에서 수동으로 경로를 체크할 필요 없이 정확한 경로 매칭이 가능합니다
5. `http.ListenAndServe`는 블로킹 함수이므로 `log.Fatal`로 에러 처리를 함께 합니다

