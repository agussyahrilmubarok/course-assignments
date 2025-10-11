package com.example.point.service.v2;

import com.example.point.domain.Point;
import com.example.point.domain.Point.PointType;
import com.example.point.domain.PointBalance;
import com.example.point.model.PointDTO;
import com.example.point.repos.PointBalanceRepository;
import com.example.point.repos.PointRepository;
import com.example.point.service.PointService;
import com.example.point.utils.UserIdInterceptor;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RLock;
import org.redisson.api.RMap;
import org.redisson.api.RedissonClient;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Service("PointServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class PointServiceImplV2 implements PointService {

    private static final String POINT_BALANCE_MAP = "point:balance";
    private static final String POINT_BALANCE_LOCK_PREFIX = "point:balance:lock:";
    private static final long LOCK_WAIT_TIME = 5L;
    private static final long LOCK_LEASE_TIME = 10L;

    private final PointBalanceRepository pointBalanceRepository;
    private final PointRepository pointRepository;
    private final RedissonClient redissonClient;

    @Override
    @Transactional
    public Point earn(PointDTO.EarnRequest request) {
        String userId = getCurrentUserId();

        if (request.getAmount() == null || request.getAmount() <= 0) {
            throw new IllegalArgumentException("Invalid amount");
        }

        RLock lock = redissonClient.getLock(POINT_BALANCE_LOCK_PREFIX + userId);
        try {
            boolean locked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!locked) throw new IllegalStateException("Failed to acquire lock for user: " + userId);

            PointBalance pointBalance = this.getOrCreatePointBalance(userId);
            pointBalance.addBalance(request.getAmount());
            pointBalance = pointBalanceRepository.save(pointBalance);

            updateCacheBalance(userId, pointBalance.getBalance());

            log.debug("User {} earned points. New balance: {}", userId, pointBalance.getBalance());
            return createAndSavePoint(userId, request.getAmount(), request.getDescription(), PointType.EARNED, pointBalance);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Lock acquisition was interrupted", e);
        } finally {
            releaseLock(lock);
        }
    }

    @Override
    @Transactional
    public Point use(PointDTO.UseRequest request) {
        String userId = getCurrentUserId();

        if (request.getAmount() == null || request.getAmount() <= 0) {
            throw new IllegalArgumentException("Invalid amount");
        }

        RLock lock = redissonClient.getLock(POINT_BALANCE_LOCK_PREFIX + userId);
        try {
            boolean locked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!locked) throw new IllegalStateException("Failed to acquire lock for user: " + userId);

            PointBalance pointBalance = this.getOrCreatePointBalance(userId);
            if (pointBalance.getBalance() < request.getAmount()) {
                log.error("Insufficient balance for user {}. Balance: {}, Required: {}", userId, pointBalance.getBalance(), request.getAmount());
                throw new IllegalArgumentException("Insufficient balance");
            }

            pointBalance.subtractBalance(request.getAmount());
            pointBalance = pointBalanceRepository.save(pointBalance);

            updateCacheBalance(userId, pointBalance.getBalance());

            log.debug("User {} used points. New balance: {}", userId, pointBalance.getBalance());
            return createAndSavePoint(userId, request.getAmount(), request.getDescription(), PointType.USED, pointBalance);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Lock acquisition was interrupted", e);
        } finally {
            releaseLock(lock);
        }
    }

    @Override
    @Transactional
    public Point cancel(PointDTO.CancelRequest request) {
        String userId = getCurrentUserId();

        Point originalPoint = pointRepository.findById(request.getPointId())
                .orElseThrow(() -> {
                    log.error("Point not found: {}", request.getPointId());
                    return new IllegalArgumentException("Point not found");
                });

        if (!originalPoint.getUserId().equals(userId)) {
            throw new IllegalArgumentException("Unauthorized to cancel this point");
        }

        if (originalPoint.getType() == PointType.CANCELED) {
            log.warn("Point already canceled: {}", request.getPointId());
            throw new IllegalArgumentException("Point already canceled");
        }

        RLock lock = redissonClient.getLock(POINT_BALANCE_LOCK_PREFIX + userId);
        try {
            boolean locked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!locked) throw new IllegalStateException("Failed to acquire lock for user: " + userId);

            PointBalance pointBalance = this.getOrCreatePointBalance(userId);
            long adjustmentAmount = originalPoint.getAmount();

            if (originalPoint.getType() == PointType.EARNED) {
                if (pointBalance.getBalance() < adjustmentAmount) {
                    throw new IllegalArgumentException("Insufficient balance to cancel earned point");
                }
                pointBalance.subtractBalance(adjustmentAmount);
            } else if (originalPoint.getType() == PointType.USED) {
                pointBalance.addBalance(adjustmentAmount);
            } else {
                throw new IllegalArgumentException("Invalid point type for cancellation");
            }

            pointBalance = pointBalanceRepository.save(pointBalance);
            updateCacheBalance(userId, pointBalance.getBalance());

            log.debug("User {} cancel points. New balance: {}", userId, pointBalance.getBalance());
            return createAndSavePoint(userId, adjustmentAmount, request.getDescription(), PointType.CANCELED, pointBalance);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Lock acquisition was interrupted", e);
        } finally {
            releaseLock(lock);
        }
    }

    @Override
    @Transactional(readOnly = true)
    public Long getBalance() {
        String userId = getCurrentUserId();
        return getBalanceInCacheOrDb(userId);
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Point> getHistory(Pageable pageable) {
        String userId = getCurrentUserId();
        return pointRepository.findByUserIdOrderByCreatedAtDesc(userId, pageable);
    }

    private String getCurrentUserId() {
        String userId = UserIdInterceptor.getCurrentUserId();
        log.debug("Current user ID: {}", userId);
        return userId;
    }

    private PointBalance getOrCreatePointBalance(String userId) {
        return pointBalanceRepository.findByUserId(userId)
                .orElseGet(() -> {
                    PointBalance balance = new PointBalance();
                    balance.setId(UUID.randomUUID().toString());
                    balance.setUserId(userId);
                    balance.setBalance(0L);
                    return pointBalanceRepository.save(balance); // perbaikan
                });
    }

    private Point createAndSavePoint(String userId, Long amount, String description,
                                     PointType pointType, PointBalance pointBalance) {
        Point point = new Point();
        point.setId(UUID.randomUUID().toString());
        point.setUserId(userId);
        point.setAmount(amount);
        point.setDescription(description);
        point.setType(pointType);
        point.setBalanceSnapshot(pointBalance.getBalance());
        point.setPointBalance(pointBalance);
        return pointRepository.save(point);
    }

    private void updateCacheBalance(String userId, Long balance) {
        RMap<String, Long> balanceMap = redissonClient.getMap(POINT_BALANCE_MAP);
        balanceMap.fastPut(userId, balance);
    }

    private Long getBalanceInCacheOrDb(String userId) {
        RMap<String, Long> balanceMap = redissonClient.getMap(POINT_BALANCE_MAP);
        Long cachedBalance = balanceMap.get(userId);
        if (cachedBalance != null) return cachedBalance;

        Long balance = pointBalanceRepository.findByUserId(userId)
                .map(PointBalance::getBalance)
                .orElse(0L);
        balanceMap.fastPut(userId, balance);
        return balance;
    }

    private void releaseLock(RLock lock) {
        if (lock != null && lock.isHeldByCurrentThread()) {
            lock.unlock();
        }
    }
}
