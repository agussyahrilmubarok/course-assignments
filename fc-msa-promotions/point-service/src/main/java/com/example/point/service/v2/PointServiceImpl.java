package com.example.point.service.v2;

import com.example.point.domain.Point;
import com.example.point.domain.PointBalance;
import com.example.point.domain.PointType;
import com.example.point.model.PointDTO;
import com.example.point.repos.PointBalanceRepository;
import com.example.point.repos.PointRepository;
import com.example.point.service.v1.PointService;
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

/**
 * V2 Point Service implementation with Redis-based locking and caching.
 */
@Service("PointServiceImplV2")
@Slf4j
@RequiredArgsConstructor
public class PointServiceImpl implements PointService {

    private static final String POINT_BALANCE_MAP = "point:balance";
    private static final String POINT_LOCK_PREFIX = "point:lock:";
    private static final long LOCK_WAIT_TIME = 3L;
    private static final long LOCK_LEASE_TIME = 3L;

    private final PointRepository pointRepository;
    private final PointBalanceRepository pointBalanceRepository;
    private final RedissonClient redissonClient;

    /**
     * Earn points and update both DB and Redis.
     */
    @Override
    @Transactional
    public Point earn(PointDTO.EarnRequest request) {
        Long userId = request.getUserId();
        Long amount = request.getAmount();
        String description = request.getDescription();

        log.info("Earning points: userId={}, amount={}", userId, amount);

        RLock lock = redissonClient.getLock(POINT_LOCK_PREFIX + userId);
        try {
            if (!lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS)) {
                throw new IllegalStateException("Could not acquire lock for user: " + userId);
            }

            Long currentBalance = getBalanceFromCache(userId);
            if (currentBalance == null) {
                currentBalance = getBalanceFromDB(userId);
                updateBalanceCache(userId, currentBalance);
            }

            PointBalance pointBalance = pointBalanceRepository.findByUserId(userId)
                    .orElseGet(() -> {
                        PointBalance pb = new PointBalance();
                        pb.setId(UUID.randomUUID().toString());
                        pb.setUserId(userId);
                        pb.setBalance(0L);
                        return pb;
                    });

            pointBalance.addBalance(amount);
            pointBalance = pointBalanceRepository.save(pointBalance);
            updateBalanceCache(userId, pointBalance.getBalance());

            Point point = new Point();
            point.setId(UUID.randomUUID().toString());
            point.setUserId(userId);
            point.setAmount(amount);
            point.setDescription(description);
            point.setType(PointType.EARNED);
            point.setBalanceSnapshot(pointBalance.getBalance());
            point.setPointBalance(pointBalance);

            return pointRepository.save(point);

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Interrupted while acquiring lock", e);
        } finally {
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
            }
        }
    }

    /**
     * Use points and update both DB and Redis.
     */
    @Override
    @Transactional
    public Point use(PointDTO.UseRequest request) {
        Long userId = request.getUserId();
        Long amount = request.getAmount();
        String description = request.getDescription();

        log.info("Using points: userId={}, amount={}", userId, amount);

        RLock lock = redissonClient.getLock(POINT_LOCK_PREFIX + userId);
        try {
            if (!lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS)) {
                throw new IllegalStateException("Could not acquire lock for user: " + userId);
            }

            Long currentBalance = getBalanceFromCache(userId);
            if (currentBalance == null) {
                currentBalance = getBalanceFromDB(userId);
                updateBalanceCache(userId, currentBalance);
            }

            if (currentBalance < amount) {
                throw new IllegalArgumentException("Insufficient balance for user: " + userId);
            }

            PointBalance pointBalance = pointBalanceRepository.findByUserId(userId)
                    .orElseThrow(() -> new IllegalArgumentException("User not found"));

            pointBalance.subtractBalance(amount);
            pointBalance = pointBalanceRepository.save(pointBalance);
            updateBalanceCache(userId, pointBalance.getBalance());

            Point point = new Point();
            point.setId(UUID.randomUUID().toString());
            point.setUserId(userId);
            point.setAmount(amount);
            point.setDescription(description);
            point.setType(PointType.USED);
            point.setBalanceSnapshot(pointBalance.getBalance());
            point.setPointBalance(pointBalance);

            return pointRepository.save(point);

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Interrupted while acquiring lock", e);
        } finally {
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
            }
        }
    }

    /**
     * Cancel previously earned/used points and update DB and Redis.
     */
    @Override
    @Transactional
    public Point cancel(PointDTO.CancelRequest request) {
        Long pointId = request.getPointId();
        String description = "Cancel pointId: " + pointId;

        Point originalPoint = pointRepository.findById(pointId)
                .orElseThrow(() -> new IllegalArgumentException("Point not found"));

        Long userId = originalPoint.getUserId();

        log.info("Cancelling point: pointId={}, userId={}", pointId, userId);

        if (originalPoint.getType() == PointType.CANCELED) {
            throw new IllegalArgumentException("Point already canceled");
        }

        RLock lock = redissonClient.getLock(POINT_LOCK_PREFIX + userId);
        try {
            if (!lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS)) {
                throw new IllegalStateException("Could not acquire lock for user: " + userId);
            }

            PointBalance pointBalance = originalPoint.getPointBalance();
            if (originalPoint.getType() == PointType.EARNED) {
                if (pointBalance.getBalance() < originalPoint.getAmount()) {
                    throw new IllegalArgumentException("Cannot cancel: insufficient balance");
                }
                pointBalance.subtractBalance(originalPoint.getAmount());
            } else if (originalPoint.getType() == PointType.USED) {
                pointBalance.addBalance(originalPoint.getAmount());
            } else {
                throw new IllegalArgumentException("Invalid point type for cancellation");
            }

            pointBalance = pointBalanceRepository.save(pointBalance);
            updateBalanceCache(userId, pointBalance.getBalance());

            Point cancelPoint = new Point();
            cancelPoint.setId(UUID.randomUUID().toString());
            cancelPoint.setUserId(userId);
            cancelPoint.setAmount(originalPoint.getAmount());
            cancelPoint.setDescription(description);
            cancelPoint.setType(PointType.CANCELED);
            cancelPoint.setBalanceSnapshot(pointBalance.getBalance());
            cancelPoint.setPointBalance(pointBalance);

            return pointRepository.save(cancelPoint);

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Interrupted while acquiring lock", e);
        } finally {
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
            }
        }
    }

    /**
     * Get user balance from cache or DB.
     */
    @Override
    @Transactional(readOnly = true)
    public Long getBalance() {
        log.warn("getBalance() called without userId. This is not supported in V2.");
        throw new UnsupportedOperationException("Please provide userId to get balance in V2.");
    }

    /**
     * Retrieve point history.
     * Not implemented here.
     */
    @Override
    public Page<Point> getHistory(Pageable pageable) {
        log.warn("getHistory() without userId is not supported.");
        throw new UnsupportedOperationException("getHistory() without userId is not supported in V2.");
    }

    // Helper methods

    private Long getBalanceFromCache(Long userId) {
        RMap<String, Long> map = redissonClient.getMap(POINT_BALANCE_MAP);
        return map.get(String.valueOf(userId));
    }

    private void updateBalanceCache(Long userId, Long balance) {
        RMap<String, Long> map = redissonClient.getMap(POINT_BALANCE_MAP);
        map.fastPut(String.valueOf(userId), balance);
    }

    @Transactional(readOnly = true)
    private Long getBalanceFromDB(Long userId) {
        return pointBalanceRepository.findByUserId(userId)
                .map(PointBalance::getBalance)
                .orElse(0L);
    }
}
