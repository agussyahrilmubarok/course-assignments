package com.example.coupon.utils;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

@Component
public class UserIdInterceptor implements HandlerInterceptor {

    private static final String USER_ID_HEADER = "X-USER-ID";
    private static final ThreadLocal<String> currentUserId = new ThreadLocal<>();

    public static String getCurrentUserId() {
        String userId = currentUserId.get();
        if (userId == null) {
            throw new IllegalStateException("User ID not found in current context");
        }
        return userId;
    }

    @Override
    public boolean preHandle(HttpServletRequest request,
                             HttpServletResponse response,
                             Object handler) throws Exception {
        String userIdStr = request.getHeader(USER_ID_HEADER);
        if (userIdStr == null || userIdStr.isEmpty()) {
            throw new IllegalStateException("X-USER-ID header is required");
        }

        try {
            currentUserId.set(userIdStr);
            return true;
        } catch (NumberFormatException e) {
            throw new IllegalStateException("Invalid X-USER-ID format");
        }
    }

    @Override
    public void afterCompletion(HttpServletRequest request,
                                HttpServletResponse response,
                                Object handler,
                                Exception ex) throws Exception {
        currentUserId.remove();
    }
}
