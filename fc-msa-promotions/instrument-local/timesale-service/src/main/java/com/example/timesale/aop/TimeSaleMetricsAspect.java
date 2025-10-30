package com.example.timesale.aop;

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
public class TimeSaleMetricsAspect {

    private final MeterRegistry registry;

    @Around("@annotation(TimeSaleMetered)")
    public Object measureTimeSaleOperation(ProceedingJoinPoint joinPoint) throws Throwable {
        Timer.Sample sample = Timer.start();
        String version = extractVersion(joinPoint);
        String operation = extractOperation(joinPoint);

        try {
            Object result = joinPoint.proceed();

            // Time sale processing success metric
            Counter.builder("timesale.operation.success")
                    .tag("version", version)
                    .tag("operation", operation)
                    .register(registry)
                    .increment();

            sample.stop(Timer.builder("timesale.operation.duration")
                    .tag("version", version)
                    .tag("operation", operation)
                    .register(registry));

            return result;
        } catch (Exception e) {
            // Time sale processing failure metric
            Counter.builder("timesale.operation.failure")
                    .tag("version", version)
                    .tag("operation", operation)
                    .tag("error", e.getClass().getSimpleName())
                    .register(registry)
                    .increment();
            throw e;
        }
    }

    private String extractVersion(ProceedingJoinPoint joinPoint) {
        TimeSaleMetered annotation = ((MethodSignature) joinPoint.getSignature())
                .getMethod()
                .getAnnotation(TimeSaleMetered.class);
        return annotation.version();
    }

    private String extractOperation(ProceedingJoinPoint joinPoint) {
        return joinPoint.getSignature().getName();
    }
}