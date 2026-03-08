package com.example.coupon.aop;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import lombok.RequiredArgsConstructor;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.stereotype.Component;

@Aspect
@Component
@RequiredArgsConstructor
public class CouponMetricsAspect {

    private final MeterRegistry registry;

    @Around("@annotation(CouponMetered)")
    public Object measureCouponOperation(ProceedingJoinPoint joinPoint) throws Throwable {
        Timer.Sample sample = Timer.start();
        String version = extractVersion(joinPoint);
        String operation = extractOperation(joinPoint);

        try {
            Object result = joinPoint.proceed();

            // Metric for successful coupon issuance
            Counter.builder("coupon.operation.success")
                    .tag("version", version)
                    .tag("operation", operation)
                    .register(registry)
                    .increment();

            sample.stop(Timer.builder("coupon.operation.duration")
                    .tag("version", version)
                    .tag("operation", operation)
                    .register(registry));

            return result;
        } catch (Exception e) {
            // Metric for failed coupon issuance
            Counter.builder("coupon.operation.failure")
                    .tag("version", version)
                    .tag("operation", operation)
                    .tag("error", e.getClass().getSimpleName())
                    .register(registry)
                    .increment();
            throw e;
        }
    }

    private String extractVersion(ProceedingJoinPoint joinPoint) {
        CouponMetered annotation = ((MethodSignature) joinPoint.getSignature())
                .getMethod()
                .getAnnotation(CouponMetered.class);
        return annotation.version();
    }

    private String extractOperation(ProceedingJoinPoint joinPoint) {
        return joinPoint.getSignature().getName();
    }
}