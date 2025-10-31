@echo off
echo ========================================
echo Building all microservices images...
echo ========================================

docker build --no-cache -f app/cart/Dockerfile -t gomall-cart:latest .
docker build --no-cache -f app/checkout/Dockerfile -t gomall-checkout:latest .
docker build --no-cache -f app/email/Dockerfile -t gomall-email:latest .
docker build --no-cache -f app/frontend/Dockerfile -t gomall-frontend:latest .
docker build --no-cache -f app/order/Dockerfile -t gomall-order:latest .
docker build --no-cache -f app/payment/Dockerfile -t gomall-payment:latest .
docker build --no-cache -f app/product/Dockerfile -t gomall-product:latest .
docker build --no-cache -f app/user/Dockerfile -t gomall-user:latest .

echo.
echo ========================================
echo All images built successfully!
echo ========================================
docker images | findstr gomall
