@echo off
echo ========================================
echo Starting all microservices...
echo ========================================

REM Create network if not exists
docker network create gomall 2>nul

REM Remove old containers if they exist
docker rm -f cart checkout email frontend order payment product user 2>nul

REM Start all microservices
echo Starting cart service...
docker run -d --name cart --network gomall -p 8883:8883 --env-file ./app/cart/.env gomall-cart:latest

echo Starting checkout service...
docker run -d --name checkout --network gomall -p 8884:8884 --env-file ./app/checkout/.env gomall-checkout:latest

echo Starting email service...
docker run -d --name email --network gomall -p 8885:8885 --env-file ./app/email/.env gomall-email:latest

echo Starting frontend service...
docker run -d --name frontend --network gomall -p 8080:8080 --env-file ./app/frontend/.env gomall-frontend:latest

echo Starting order service...
docker run -d --name order --network gomall -p 8886:8886 --env-file ./app/order/.env gomall-order:latest

echo Starting payment service...
docker run -d --name payment --network gomall -p 8887:8887 --env-file ./app/payment/.env gomall-payment:latest

echo Starting product service...
docker run -d --name product --network gomall -p 8888:8888 --env-file ./app/product/.env gomall-product:latest

echo Starting user service...
docker run -d --name user --network gomall -p 8889:8889 --env-file ./app/user/.env gomall-user:latest

echo.
echo ========================================
echo All services started!
echo ========================================
docker ps | findstr gomall
echo.
echo Frontend is available at: http://localhost:8080
