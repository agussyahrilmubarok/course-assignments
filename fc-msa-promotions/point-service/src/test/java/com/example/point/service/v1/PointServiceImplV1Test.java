package com.example.point.service.v1;

import com.example.point.domain.Point;
import com.example.point.domain.Point.PointType;
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
import org.mockito.Mockito;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class PointServiceImplV1Test {

    private static final String TEST_USER_ID = "USER_1";
    private static final String TEST_POINT_ID = "POINT_1";

    @InjectMocks
    private PointServiceImplV1 pointService;

    @Mock
    private PointBalanceRepository pointBalanceRepository;

    @Mock
    private PointRepository pointRepository;

    private PointBalance pointBalance;
    private Point point;

    @BeforeEach
    void setUp() {
        pointBalance = new PointBalance();
        pointBalance.setId(UUID.randomUUID().toString());
        pointBalance.setUserId(TEST_USER_ID);
        pointBalance.setBalance(5000L);

        point = new Point();
        point.setId(TEST_POINT_ID);
        point.setUserId(TEST_USER_ID);
        point.setAmount(1000L);
        point.setType(PointType.EARNED);
        point.setBalanceSnapshot(6000L);
        point.setPointBalance(pointBalance);
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
            when(pointBalanceRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

            Point result = pointService.earn(request);

            assertThat(result)
                    .isNotNull()
                    .extracting(Point::getUserId, Point::getAmount, Point::getType)
                    .containsExactly(TEST_USER_ID, 1000L, PointType.EARNED);

            assertThat(result.getBalanceSnapshot()).isEqualTo(pointBalance.getBalance());

            verify(pointBalanceRepository).save(any());
            verify(pointRepository).save(any());
        }
    }

    @Test
    void shouldCreateNewBalanceWhenNoneExistsAndEarnPoint() {
        PointDTO.EarnRequest request = PointDTO.EarnRequest.builder()
                .amount(1000L)
                .description("First earn")
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.empty());
            when(pointBalanceRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

            Point result = pointService.earn(request);

            assertThat(result).isNotNull();
            assertThat(result.getType()).isEqualTo(PointType.EARNED);
            assertThat(result.getAmount()).isEqualTo(1000L);

            verify(pointBalanceRepository, times(2)).save(any());
            verify(pointRepository).save(any());
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
            when(pointBalanceRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

            Point result = pointService.use(request);

            assertThat(result).isNotNull();
            assertThat(result.getType()).isEqualTo(PointType.USED);
            assertThat(result.getAmount()).isEqualTo(1000L);
            assertThat(result.getBalanceSnapshot()).isEqualTo(pointBalance.getBalance());

            verify(pointBalanceRepository).save(any());
            verify(pointRepository).save(any());
        }
    }

    @Test
    void shouldThrowExceptionWhenNoBalanceOnUse() {
        PointDTO.UseRequest request = PointDTO.UseRequest.builder()
                .amount(1000L)
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.empty());

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.use(request));
            assertThat(ex.getMessage()).isEqualTo("Point balance not found");
        }
    }

    @Test
    void shouldThrowExceptionWhenInsufficientBalanceOnUse() {
        PointDTO.UseRequest request = PointDTO.UseRequest.builder()
                .amount(10000L)
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.use(request));
            assertThat(ex.getMessage()).isEqualTo("Insufficient balance");
        }
    }

    @Test
    void shouldCancelEarnedPointAndDeductBalance() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();

        point.setType(PointType.EARNED);
        point.setAmount(1000L);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));
            when(pointBalanceRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

            Point result = pointService.cancel(request);

            assertThat(result.getType()).isEqualTo(PointType.CANCELED);
            assertThat(result.getAmount()).isEqualTo(1000L);
        }
    }

    @Test
    void shouldThrowExceptionWhenPointAlreadyCanceled() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();
        point.setType(PointType.CANCELED);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.cancel(request));
            assertThat(ex.getMessage()).isEqualTo("Point already canceled");
        }
    }

    @Test
    void shouldCancelUsedPointAndAddBalance() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();

        point.setType(PointType.USED);
        point.setAmount(1000L);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));
            when(pointBalanceRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

            Point result = pointService.cancel(request);

            assertThat(result.getType()).isEqualTo(PointType.CANCELED);
            assertThat(result.getAmount()).isEqualTo(1000L);
        }
    }

    @Test
    void shouldReturnCorrectBalanceWhenExists() {
        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));

            Long result = pointService.getBalance();

            assertThat(result).isEqualTo(5000L);
        }
    }

    @Test
    void shouldReturnZeroBalanceWhenNoneExists() {
        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.empty());

            Long result = pointService.getBalance();

            assertThat(result).isZero();
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

    @Test
    void shouldThrowExceptionWhenPointNotFoundOnCancel() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.empty());

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.cancel(request));
            assertThat(ex.getMessage()).isEqualTo("Point not found");
        }
    }

    @Test
    void shouldThrowExceptionWhenCancelEarnedPointWithInsufficientBalance() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();

        point.setType(PointType.EARNED);
        point.setAmount(6000L);
        pointBalance.setBalance(5000L);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.cancel(request));
            assertThat(ex.getMessage()).isEqualTo("Insufficient balance to cancel earned point");
        }
    }

    @Test
    void shouldThrowExceptionWhenPointTypeInvalidOnCancel() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();

        point.setType(null);
        point.setUserId(TEST_USER_ID);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);
            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(pointBalance));

            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> pointService.cancel(request));
            assertThat(ex.getMessage()).isEqualTo("Invalid point type for cancellation");
        }
    }

    @Test
    void shouldCancelEarnedPointAndSubtractBalance() {
        PointDTO.CancelRequest request = PointDTO.CancelRequest.builder()
                .pointId(TEST_POINT_ID)
                .build();

        point.setType(PointType.EARNED);
        point.setAmount(1000L);
        point.setUserId(TEST_USER_ID);
        pointBalance.setBalance(5000L);

        PointBalance spyBalance = Mockito.spy(pointBalance);

        try (MockedStatic<UserIdInterceptor> mockUser = mockStatic(UserIdInterceptor.class)) {
            mockUser.when(UserIdInterceptor::getCurrentUserId).thenReturn(TEST_USER_ID);

            when(pointRepository.findById(TEST_POINT_ID)).thenReturn(Optional.of(point));
            when(pointBalanceRepository.findByUserId(TEST_USER_ID)).thenReturn(Optional.of(spyBalance));
            when(pointBalanceRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));
            when(pointRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

            Point result = pointService.cancel(request);

            assertNotNull(result);
            assertEquals(PointType.CANCELED, result.getType());
            verify(spyBalance).subtractBalance(1000L);
        }
    }
}
