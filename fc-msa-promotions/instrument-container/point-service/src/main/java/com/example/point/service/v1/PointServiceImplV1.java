package com.example.point.service.v1;

import com.example.point.aop.PointMetered;
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
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Isolation;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service("PointServiceImplV1")
@Slf4j
@RequiredArgsConstructor
public class PointServiceImplV1 implements PointService {

    private final PointBalanceRepository pointBalanceRepository;
    private final PointRepository pointRepository;

    @Override
    @Transactional(isolation = Isolation.REPEATABLE_READ)
    @PointMetered(version = "v1")
    public Point earn(PointDTO.EarnRequest request) {
        validateAmount(request.getAmount());

        String userId = getCurrentUserId();
        PointBalance pointBalance = getOrCreatePointBalance(userId);

        pointBalance.addBalance(request.getAmount());
        pointBalance = pointBalanceRepository.save(pointBalance);

        log.debug("User {} earned points. New balance: {}", userId, pointBalance.getBalance());

        return createAndSavePoint(userId, request.getAmount(), request.getDescription(), PointType.EARNED, pointBalance);
    }

    @Override
    @Transactional(isolation = Isolation.REPEATABLE_READ)
    @PointMetered(version = "v1")
    public Point use(PointDTO.UseRequest request) {
        validateAmount(request.getAmount());

        String userId = getCurrentUserId();
        PointBalance pointBalance = getPointBalanceOrThrow(userId);

        if (pointBalance.getBalance() < request.getAmount()) {
            log.error("Insufficient balance for user {}. Balance: {}, Required: {}", userId, pointBalance.getBalance(), request.getAmount());
            throw new IllegalArgumentException("Insufficient balance");
        }

        pointBalance.subtractBalance(request.getAmount());
        pointBalance = pointBalanceRepository.save(pointBalance);

        log.debug("User {} used points. New balance: {}", userId, pointBalance.getBalance());

        return createAndSavePoint(userId, request.getAmount(), request.getDescription(), PointType.USED, pointBalance);
    }

    @Override
    @Transactional(isolation = Isolation.REPEATABLE_READ)
    public Point cancel(PointDTO.CancelRequest request) {
        String userId = getCurrentUserId();

        Point originalPoint = pointRepository.findById(request.getPointId())
                .orElseThrow(() -> {
                    log.error("Point not found: {}", request.getPointId());
                    return new IllegalArgumentException("Point not found");
                });

        if (!originalPoint.getUserId().equals(userId)) {
            log.error("User {} is not the owner of point {}", userId, request.getPointId());
            throw new IllegalArgumentException("Unauthorized to cancel this point");
        }

        if (originalPoint.getType() == PointType.CANCELED) {
            log.warn("Point already canceled: {}", request.getPointId());
            throw new IllegalArgumentException("Point already canceled");
        }

        PointBalance pointBalance = getPointBalanceOrThrow(userId);
        long adjustmentAmount = originalPoint.getAmount();

        if (originalPoint.getType() == PointType.EARNED) {
            if (pointBalance.getBalance() < adjustmentAmount) {
                log.error("Cannot cancel earned point due to insufficient balance. User: {}, Balance: {}, Required: {}",
                        userId, pointBalance.getBalance(), adjustmentAmount);
                throw new IllegalArgumentException("Insufficient balance to cancel earned point");
            }
            pointBalance.subtractBalance(adjustmentAmount);
        } else if (originalPoint.getType() == PointType.USED) {
            pointBalance.addBalance(adjustmentAmount);
        } else {
            log.error("Invalid point type for cancellation: {}, Type: {}", request.getPointId(), originalPoint.getType());
            throw new IllegalArgumentException("Invalid point type for cancellation");
        }

        pointBalance = pointBalanceRepository.save(pointBalance);

        return createAndSavePoint(userId, adjustmentAmount, request.getDescription(), PointType.CANCELED, pointBalance);
    }

    @Override
    @Transactional(readOnly = true)
    public Long getBalance() {
        String userId = getCurrentUserId();
        return pointBalanceRepository.findByUserId(userId)
                .map(PointBalance::getBalance)
                .orElse(0L);
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
                    log.info("Creating new point balance for user: {}", userId);
                    PointBalance balance = new PointBalance();
                    balance.setId(UUID.randomUUID().toString());
                    balance.setUserId(userId);
                    balance.setBalance(0L);
                    return pointBalanceRepository.save(balance);
                });
    }

    private PointBalance getPointBalanceOrThrow(String userId) {
        return pointBalanceRepository.findByUserId(userId)
                .orElseThrow(() -> {
                    log.error("Point balance not found for user: {}", userId);
                    return new IllegalArgumentException("Point balance not found");
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
        point = pointRepository.save(point);

        log.info("Point created: ID={}, Type={}, Amount={}, BalanceSnapshot={}, User={}",
                point.getId(), pointType, amount, point.getBalanceSnapshot(), userId);

        return point;
    }

    private void validateAmount(Long amount) {
        if (amount == null || amount <= 0) {
            throw new IllegalArgumentException("Invalid amount");
        }
    }
}
