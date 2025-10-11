package com.example.point.service.v2;

import com.example.point.domain.Point;
import com.example.point.domain.PointBalance;
import com.example.point.model.PointDTO;
import com.example.point.repos.PointBalanceRepository;
import com.example.point.repos.PointRepository;
import com.example.point.utils.UserIdInterceptor;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.junit.jupiter.MockitoExtension;
import org.redisson.api.RLock;
import org.redisson.api.RMap;
import org.redisson.api.RedissonClient;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class PointServiceImplV2Test {

    private static final String TEST_USER_ID = "USER_1";
    private static final String TEST_POINT_ID = "POINT_1";

    @InjectMocks
    private PointServiceImplV2 pointService;

    @Mock
    private PointBalanceRepository pointBalanceRepository;

    @Mock
    private PointRepository pointRepository;

    @Mock
    private RedissonClient redissonClient;

    @Mock
    private RLock lock;

    @Mock
    private RMap<Object, Object> balanceMap;

    private PointBalance pointBalance;
    private Point point;

    @BeforeEach
    void setUp() throws InterruptedException {
        pointBalance = new PointBalance();
        pointBalance.setId(UUID.randomUUID().toString());
        pointBalance.setUserId(TEST_USER_ID);
        pointBalance.setBalance(5000L);

        point = new Point();
        point.setId(TEST_POINT_ID);
        point.setUserId(TEST_USER_ID);
        point.setAmount(1000L);
        point.setType(Point.PointType.EARNED);
        point.setBalanceSnapshot(6000L);
        point.setPointBalance(pointBalance);

        lenient().when(redissonClient.getLock(anyString())).thenReturn(lock);
        lenient().when(redissonClient.getMap(anyString())).thenReturn(balanceMap);
        lenient().when(lock.tryLock(anyLong(), anyLong(), any(TimeUnit.class))).thenReturn(true);
        lenient().when(lock.isHeldByCurrentThread()).thenReturn(true);
    }

    @Test
    void shouldEarnPointAndUpdateBalance() {
        PointDTO.EarnRequest request = PointDTO.EarnRequest.builder()
                .amount(1000L)
                .description("Earned test point")
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));
            when(pointBalanceRepository.save(any())).thenAnswer(i -> i.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(i -> i.getArgument(0));

            Point result = pointService.earn(request);

            assertThat(result).isNotNull();
            assertThat(result.getUserId()).isEqualTo(TEST_USER_ID);
            assertThat(result.getType()).isEqualTo(Point.PointType.EARNED);

            verify(lock).tryLock(anyLong(), anyLong(), any());
            verify(pointBalanceRepository).save(any());
            verify(pointRepository).save(any());
            verify(balanceMap).fastPut(TEST_USER_ID, String.valueOf(pointBalance.getBalance()));
            verify(lock).unlock();
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }
    }

    @Test
    void shouldUsePointAndDeductBalance() {
        PointDTO.UseRequest request = PointDTO.UseRequest.builder()
                .amount(1000L)
                .description("Use points")
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));
            when(pointBalanceRepository.save(any())).thenAnswer(i -> i.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(i -> i.getArgument(0));

            Point result = pointService.use(request);

            assertThat(result).isNotNull();
            assertThat(result.getType()).isEqualTo(Point.PointType.USED);

            verify(pointBalanceRepository).save(any());
            verify(balanceMap).fastPut(TEST_USER_ID, String.valueOf(pointBalance.getBalance()));
            verify(lock).unlock();
        }
    }

    @Test
    void shouldThrowExceptionWhenInsufficientBalanceOnUse() {
        PointDTO.UseRequest request = PointDTO.UseRequest.builder()
                .amount(99999L)
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.use(request));
            assertThat(ex.getMessage()).isEqualTo("Insufficient balance");
        }
    }

    @Test
    void shouldCancelUsedPointAndAddBalance() {
        point.setType(Point.PointType.USED);
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder().pointId(TEST_POINT_ID).build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));
            when(pointBalanceRepository.save(any())).thenAnswer(i -> i.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(i -> i.getArgument(0));

            Point result = pointService.cancel(request);

            assertThat(result).isNotNull();
            assertThat(result.getType()).isEqualTo(Point.PointType.CANCELED);
            verify(balanceMap).fastPut(TEST_USER_ID, String.valueOf(pointBalance.getBalance()));
            verify(lock).unlock();
        }
    }

    @Test
    void shouldThrowExceptionWhenLockFails() throws InterruptedException {
        PointDTO.EarnRequest request = PointDTO.EarnRequest.builder()
                .amount(500L)
                .build();

        when(lock.tryLock(anyLong(), anyLong(), any())).thenReturn(false);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            IllegalStateException ex = assertThrows(IllegalStateException.class, () -> pointService.earn(request));
            assertThat(ex.getMessage()).contains("Failed to acquire lock");
        }
    }

    @Test
    void shouldReturnCachedBalance() {
        when(balanceMap.get(TEST_USER_ID)).thenReturn(String.valueOf(7000L));
        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            Long result = pointService.getBalance();
            assertThat(result).isEqualTo(7000L);
        }
    }

    @Test
    void shouldReturnDbBalanceWhenCacheMiss() {
        when(balanceMap.get(TEST_USER_ID)).thenReturn(null);
        when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            Long result = pointService.getBalance();
            assertThat(result).isEqualTo(5000L);
            verify(balanceMap).fastPut(TEST_USER_ID, String.valueOf(5000L));
        }
    }

    @Test
    void shouldReturnUserPointHistory() {
        Pageable pageable = PageRequest.of(0, 10);
        Page<Point> page = new PageImpl<>(List.of(point));

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findByUserIdOrderByCreatedAtDesc(TEST_USER_ID, pageable)).thenReturn(page);

            Page<Point> result = pointService.getHistory(pageable);

            assertThat(result.getContent()).hasSize(1);
            assertThat(result.getContent().getFirst().getId()).isEqualTo(TEST_POINT_ID);
        }
    }
}
