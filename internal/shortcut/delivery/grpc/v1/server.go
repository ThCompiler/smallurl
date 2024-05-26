package v1

import (
	"context"
	"net/url"
	"smallurl/internal/shortcut/delivery"
	sr "smallurl/internal/shortcut/repository"
	generatedv1 "smallurl/pkg/grpc/v1"
	"smallurl/pkg/logger"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortcutService struct {
	generatedv1.UnimplementedShortcutServer
	usc delivery.Usecase
	l   logger.Interface
}

func NewShortcutService(usc delivery.Usecase, l logger.Interface) *ShortcutService {
	return &ShortcutService{
		usc: usc,
		l:   l,
	}
}

func (s *ShortcutService) GetShortURL(_ context.Context, longURL *generatedv1.LongUrl) (*generatedv1.ShortUrl, error) {
	// Проверка валидности URL
	if _, err := url.ParseRequestURI(longURL.GetLongUrl()); err != nil {
		s.l.Info("[GRPC} Result - error %s was sent with status code %d", ErrorURLNotValid, codes.InvalidArgument)

		return nil, status.Error(codes.InvalidArgument, ErrorURLNotValid.Error())
	}

	shortURL, err := s.usc.GetShortURL(longURL.GetLongUrl())
	if err != nil {
		s.l.Error("[GRPC} - %s", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &generatedv1.ShortUrl{
		ShortUrl: shortURL,
	}, nil
}

func (s *ShortcutService) GetLongURL(_ context.Context, shortURL *generatedv1.ShortUrl) (*generatedv1.LongUrl, error) {
	longURL, err := s.usc.GetLongURL(shortURL.GetShortUrl())
	if err != nil {
		if errors.Is(err, sr.ErrorURLNotFound) {
			s.l.Info("[GRPC} Result - error %s was sent with status code %d", err, codes.NotFound)

			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.l.Error("[GRPC} - %s", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &generatedv1.LongUrl{
		LongUrl: longURL,
	}, nil
}
