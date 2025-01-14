#!/bin/sh

# Generate mocks for repository interfaces
mockery --name=TenantRepository --dir=internal/repository --output=internal/test/mockrepository --outpkg=mockrepository

# Generate mocks for service interfaces
mockery --name=Messagging --dir=internal/service/messaging --output=internal/test/mockservice --outpkg=mockservice

# Generate mocks for usecase interfaces
mockery --name=TenantUsecase --dir=internal/usecase --output=internal/test/mockusecase --outpkg=mockusecase