package jwtToken

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 비밀키 (실제론 환경 변수나 설정 파일에서 가져오는 게 안전)
var jwtSecret = []byte("jwt secret key") // 환경 변수로 바꾸는 걸 추천

func VsalidateJWT(tokenString string) error {
	// "Bearer " 접두사 제거 (gRPC 메타데이터에서 이미 제거된 상태라면 생략 가능)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// JWT 파싱 및 검증
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 서명 방식이 HMAC-SHA256인지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// 파싱 중 에러 발생 시
	if err != nil {
		return fmt.Errorf("failed to parse token: %v", err)
	}

	// 토큰이 유효한지 확인
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// 클레임 추출 (선택 사항)
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return fmt.Errorf("invalid claims format")
	}

	// 만료 시간 체크 (jwt.Parse에서 자동으로 체크되지만, 명시적으로 확인 가능)
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token has expired")
	}

	// 추가 검증 (필요 시)
	// 예: Issuer(iss)나 Audience(aud) 체크
	// if claims.Issuer != "your-app" {
	//     return fmt.Errorf("invalid issuer")
	// }

	return nil // 유효하면 nil 반환
}
