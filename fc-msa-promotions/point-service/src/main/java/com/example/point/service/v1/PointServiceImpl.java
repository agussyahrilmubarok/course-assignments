package com.example.point.service.v1;

import com.example.point.domain.Point;
import com.example.point.domain.Point.PointType;
import com.example.point.domain.PointBalance;
import com.example.point.model.PointDTO;
import com.example.point.repos.PointBalanceRepository;
import com.example.point.repos.PointRepository;
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
public class PointServiceImpl implements PointService {

    private final PointBalanceRepository pointBalanceRepository;
    private final PointRepository pointRepository;

    @Override
    @Transactional(isolation = Isolation.REPEATABLE_READ)
    public Point earn(PointDTO.EarnRequest request) {
        log.info("Earning point for userId: {}, amount: {}", request.getUserId(), request.getAmount());

        // Find or create PointBalance
        PointBalance pointBalance = pointBalanceRepository.findByUserId(request.getUserId())
                .orElseGet(() -> {
                    log.debug("Creating new PointBalance for userId: {}", request.getUserId());
                    PointBalance balance = new PointBalance();
                    balance.setId(UUID.randomUUID().toString());
                    balance.setUserId(request.getUserId());
                    balance.setBalance(0L);
                    return balance;
                });

        // Add balance
        pointBalance.addBalance(request.getAmount());
        pointBalance = pointBalanceRepository.save(pointBalance);
        log.debug("Updated balance: {}", pointBalance.getBalance());

        // Create point record
        Point point = new Point();
        point.setId(UUID.randomUUID().toString());
        point.setAmount(request.getAmount());
        point.setUserId(request.getUserId());
        point.setDescription(request.getDescription());
        point.setType(PointType.EARNED);
        point.setBalanceSnapshot(pointBalance.getBalance());
        point.setPointBalance(pointBalance);

        log.info("Earn point recorded: {}", point);
        return pointRepository.save(point);
    }

    @Override
    @Transactional
    public Point use(PointDTO.UseRequest request) {
        log.info("Using point for userId: {}, amount: {}", request.getUserId(), request.getAmount());

        // Retrieve user balance
        PointBalance pointBalance = pointBalanceRepository.findByUserId(request.getUserId())
                .orElseThrow(() -> {
                    log.error("User not found: {}", request.getUserId());
                    return new IllegalArgumentException("User not found");
                });

        if (pointBalance.getBalance() < request.getAmount()) {
            log.error("Insufficient balance. Current: {}, Required: {}", pointBalance.getBalance(), request.getAmount());
            throw new IllegalArgumentException("Insufficient balance");
        }

        // Subtract balance
        pointBalance.subtractBalance(request.getAmount());
        pointBalance = pointBalanceRepository.save(pointBalance);
        log.debug("New balance after use: {}", pointBalance.getBalance());

        // Create point record
        Point point = new Point();
        point.setId(UUID.randomUUID().toString());
        point.setUserId(request.getUserId());
        point.setAmount(request.getAmount());
        point.setDescription(request.getDescription());
        point.setType(PointType.USED);
        point.setBalanceSnapshot(pointBalance.getBalance());
        point.setPointBalance(pointBalance);

        log.info("Use point recorded: {}", point);
        return pointRepository.save(point);
    }

    @Override
    @Transactional(isolation = Isolation.REPEATABLE_READ)
    public Point cancel(PointDTO.CancelRequest request) {
        log.info("Cancelling pointId: {}", request.getPointId());

        // Find original point
        Point originalPoint = pointRepository.findById(request.getPointId())
                .orElseThrow(() -> {
                    log.error("Point not found: {}", request.getPointId());
                    return new IllegalArgumentException("Point not found");
                });

        if (originalPoint.getType() == PointType.CANCELED) {
            log.warn("Attempted to cancel already canceled point: {}", request.getPointId());
            throw new IllegalArgumentException("Point already canceled");
        }

        // Retrieve balance
        PointBalance pointBalance = pointBalanceRepository.findByUserId(originalPoint.getUserId())
                .orElseThrow(() -> {
                    log.error("User not found for cancellation: {}", originalPoint.getUserId());
                    return new IllegalArgumentException("User not found");
                });

        Long newBalance;

        // Adjust balance based on original point type
        if (originalPoint.getType() == PointType.EARNED) {
            if (pointBalance.getBalance() < originalPoint.getAmount()) {
                log.error("Cannot cancel earned point: insufficient balance");
                throw new IllegalArgumentException("Cannot cancel earned point: insufficient balance");
            }
            newBalance = pointBalance.getBalance() - originalPoint.getAmount();
        } else if (originalPoint.getType() == PointType.USED) {
            newBalance = pointBalance.getBalance() + originalPoint.getAmount();
        } else {
            log.error("Invalid point type for cancellation: {}", originalPoint.getType());
            throw new IllegalArgumentException("Invalid point type for cancellation");
        }

        pointBalance.setBalance(newBalance);
        pointBalance = pointBalanceRepository.save(pointBalance);
        log.debug("Balance after cancellation: {}", pointBalance.getBalance());

        // Create cancel point
        Point cancelPoint = new Point();
        cancelPoint.setId(UUID.randomUUID().toString());
        cancelPoint.setUserId(originalPoint.getUserId());
        cancelPoint.setAmount(originalPoint.getAmount());
        cancelPoint.setDescription("Cancel pointId: " + originalPoint.getId());
        cancelPoint.setType(PointType.CANCELED);
        cancelPoint.setBalanceSnapshot(pointBalance.getBalance());
        cancelPoint.setPointBalance(pointBalance);

        log.info("Cancel point recorded: {}", cancelPoint);
        return pointRepository.save(cancelPoint);
    }

    @Override
    @Transactional(readOnly = true)
    public Long getBalance(String userId) {
        log.info("Getting balance for userId: {}", userId);
        return pointBalanceRepository.findByUserId(userId)
                .map(PointBalance::getBalance)
                .orElse(0L);
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Point> getHistory(String userId, Pageable pageable) {
        log.info("Getting point history for userId: {}", userId);
        return pointRepository.findByUserIdOrderByCreatedAtDesc(userId, pageable);
    }
}
