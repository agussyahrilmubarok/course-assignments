package com.example.timesale.service.v3;

import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.domain.TimeSaleOrder;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.redisson.api.RAtomicLong;
import org.redisson.api.RBucket;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class TimeSaleServiceImplV3Test {

    @InjectMocks
    private TimeSaleServiceImplV3 timeSaleService;

    @Mock
    private TimeSaleRepository timeSaleRepository;
    @Mock
    private ProductRepository productRepository;
    @Mock
    private TimeSaleOrderRepository timeSaleOrderRepository;
    @Mock
    private RedissonClient redissonClient;
    @Mock
    private ObjectMapper objectMapper;
    @Mock
    private RLock rLock;
    @Mock
    private RBucket<Object> rBucket;
    @Mock
    private RAtomicLong rAtomicLong;

    @Test
    void testCreate_whenValidRequest_shouldReturnSavedTimeSale() throws JsonProcessingException {
        TimeSaleDTO.ProductRequest productReq = TimeSaleDTO.ProductRequest.builder()
                .name("Prod A")
                .price(1000L)
                .build();
        TimeSaleDTO.CreateRequest request = TimeSaleDTO.CreateRequest.builder()
                .product(productReq)
                .quantity(10L)
                .discountPrice(800L)
                .startAt(LocalDateTime.now().minusMinutes(1))
                .endAt(LocalDateTime.now().plusHours(1))
                .build();

        Product savedProduct = new Product();
        savedProduct.setId("prod-1");
        savedProduct.setName("Prod A");
        savedProduct.setPrice(1000L);
        when(productRepository.save(any(Product.class))).thenReturn(savedProduct);

        TimeSale savedTimeSale = new TimeSale();
        savedTimeSale.setId("ts-1");
        savedTimeSale.setQuantity(10L);
        savedTimeSale.setRemainingQuantity(10L);
        savedTimeSale.setDiscountPrice(800L);
        savedTimeSale.setStartAt(request.getStartAt());
        savedTimeSale.setEndAt(request.getEndAt());
        savedTimeSale.setProduct(savedProduct);
        savedTimeSale.setStatus(TimeSale.Status.ACTIVE);
        when(timeSaleRepository.save(any(TimeSale.class))).thenReturn(savedTimeSale);

        when(objectMapper.writeValueAsString(any(TimeSale.class))).thenReturn("{}");
        when(redissonClient.getBucket("time-sale:" + savedTimeSale.getId())).thenReturn(rBucket);

        TimeSale result = timeSaleService.create(request);

        assertNotNull(result);
        assertEquals(10L, result.getQuantity());
        verify(productRepository).save(any(Product.class));
        verify(timeSaleRepository).save(any(TimeSale.class));
        verify(objectMapper).writeValueAsString(savedTimeSale);
        verify(rBucket).set("{}");
    }

    @Test
    void testCreate_whenStartTimeAfterEndTime_shouldThrowException() {
        TimeSaleDTO.CreateRequest request = TimeSaleDTO.CreateRequest.builder()
                .product(TimeSaleDTO.ProductRequest.builder().name("P1").price(1000L).build())
                .quantity(5L)
                .discountPrice(900L)
                .startAt(LocalDateTime.now().plusHours(2))
                .endAt(LocalDateTime.now().plusHours(1))
                .build();

        Exception ex = assertThrows(IllegalArgumentException.class,
                () -> timeSaleService.create(request));
        assertEquals("Start time must be before end time", ex.getMessage());
    }

    @Test
    void testFindById_whenExistsInCache_shouldReturnTimeSale() throws JsonProcessingException {
        String timeSaleId = "ts-123";
        TimeSale cached = new TimeSale();
        cached.setId(timeSaleId);

        when(redissonClient.getBucket("time-sale:" + timeSaleId)).thenReturn(rBucket);
        when(rBucket.get()).thenReturn("{}");
        when(objectMapper.readValue("{}", TimeSale.class)).thenReturn(cached);

        TimeSale result = timeSaleService.findById(timeSaleId);

        assertNotNull(result);
        assertEquals(timeSaleId, result.getId());
        verify(timeSaleRepository, never()).findById(anyString());
    }

    @Test
    void testFindById_whenCacheMiss_shouldLoadFromDbAndSetCache() throws JsonProcessingException {
        String timeSaleId = "ts-123";
        TimeSale fromDb = new TimeSale();
        fromDb.setId(timeSaleId);

        when(redissonClient.getBucket("time-sale:" + timeSaleId)).thenReturn(rBucket);
        when(rBucket.get()).thenReturn(null);
        when(timeSaleRepository.findById(timeSaleId)).thenReturn(Optional.of(fromDb));
        when(objectMapper.writeValueAsString(fromDb)).thenReturn("{}");

        TimeSale result = timeSaleService.findById(timeSaleId);

        assertNotNull(result);
        assertEquals(timeSaleId, result.getId());
        verify(timeSaleRepository).findById(timeSaleId);
        verify(rBucket).set("{}");
    }

    @Test
    void testFindById_whenNotFound_shouldThrowException() {
        String timeSaleId = "invalid";

        when(redissonClient.getBucket("time-sale:" + timeSaleId)).thenReturn(rBucket);
        when(rBucket.get()).thenReturn(null);
        when(timeSaleRepository.findById(timeSaleId)).thenReturn(Optional.empty());

        Exception ex = assertThrows(IllegalArgumentException.class,
                () -> timeSaleService.findById(timeSaleId));
        assertEquals("TimeSale not found", ex.getMessage());
    }

    @Test
    void testFindAllOngoing_whenCalled_shouldReturnPageOfTimeSales() {
        Pageable pageable = PageRequest.of(0, 10);
        TimeSale ts = new TimeSale();
        ts.setId("ts-1");
        Page<TimeSale> expected = new PageImpl<>(List.of(ts));

        when(timeSaleRepository.findAllByStartAtBeforeAndEndAtAfterAndStatus(any(LocalDateTime.class),
                eq(TimeSale.Status.ACTIVE), eq(pageable))).thenReturn(expected);

        Page<TimeSale> result = timeSaleService.findAllOngoing(pageable);

        assertNotNull(result);
        assertEquals(1, result.getContent().size());
    }

    @Test
    void testPurchase_whenValidRequest_shouldReturnUpdatedTimeSale() throws InterruptedException, JsonProcessingException {
        String userId = "user-1";
        String timeSaleId = "ts-1";
        TimeSaleDTO.PurchaseRequest request = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(timeSaleId)
                .quantity(2L)
                .build();

        Product prod = new Product();
        prod.setId("prod-1");
        prod.setName("Product 1");
        prod.setPrice(1000L);

        TimeSale timeSale = new TimeSale();
        timeSale.setId(timeSaleId);
        timeSale.setQuantity(10L);
        timeSale.setRemainingQuantity(10L);
        timeSale.setDiscountPrice(800L);
        timeSale.setStartAt(LocalDateTime.now().minusMinutes(10));
        timeSale.setEndAt(LocalDateTime.now().plusHours(1));
        timeSale.setStatus(TimeSale.Status.ACTIVE);
        timeSale.setProduct(prod);

        when(redissonClient.getLock("time-sale-lock:" + timeSaleId)).thenReturn(rLock);
        when(rLock.tryLock(3L, 3L, TimeUnit.SECONDS)).thenReturn(true);

        when(redissonClient.getBucket("time-sale:" + timeSaleId)).thenReturn(rBucket);
        when(rBucket.get()).thenReturn(null);
        when(timeSaleRepository.findById(timeSaleId)).thenReturn(Optional.of(timeSale));
        when(timeSaleRepository.save(any(TimeSale.class))).thenAnswer(inv -> inv.getArgument(0));
        when(timeSaleOrderRepository.save(any(TimeSaleOrder.class))).thenReturn(new TimeSaleOrder());
        when(objectMapper.writeValueAsString(any(TimeSale.class))).thenReturn("{}");

        TimeSale result = timeSaleService.purchase(request, userId);

        assertNotNull(result);
        assertEquals(8L, result.getRemainingQuantity());
        verify(timeSaleOrderRepository).save(any(TimeSaleOrder.class));
        verify(rLock).unlock();
        verify(rBucket, times(2)).set("{}");
    }

    @Test
    void testPurchase_whenLockNotAcquired_shouldThrowException() throws InterruptedException {
        String timeSaleId = "ts-1";
        TimeSaleDTO.PurchaseRequest request = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(timeSaleId)
                .quantity(1L)
                .build();

        when(redissonClient.getLock("time-sale-lock:" + timeSaleId)).thenReturn(rLock);
        when(rLock.tryLock(3L, 3L, TimeUnit.SECONDS)).thenReturn(false);

        IllegalStateException ex = assertThrows(IllegalStateException.class,
                () -> timeSaleService.purchase(request, "user-123"));
        assertEquals("Could not acquire lock for time sale purchase.", ex.getMessage());
        verify(rLock, never()).unlock();
    }

    @Test
    void testPurchase_whenInterrupted_shouldThrowException() throws InterruptedException {
        String timeSaleId = "ts-1";
        TimeSaleDTO.PurchaseRequest request = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(timeSaleId)
                .quantity(1L)
                .build();

        when(redissonClient.getLock("time-sale-lock:" + timeSaleId)).thenReturn(rLock);
        when(rLock.tryLock(3L, 3L, TimeUnit.SECONDS)).thenThrow(new InterruptedException("Interrupted"));

        IllegalStateException ex = assertThrows(IllegalStateException.class,
                () -> timeSaleService.purchase(request, "user-123"));
        assertEquals("Thread was interrupted while waiting for lock.", ex.getMessage());
        verify(rLock, never()).unlock();
    }

    @Test
    void testPurchase_whenNotEnoughStock_shouldThrowException() throws InterruptedException {
        String userId = "user-1";
        String timeSaleId = "ts-1";
        TimeSaleDTO.PurchaseRequest request = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(timeSaleId)
                .quantity(20L)
                .build();

        TimeSale timeSale = new TimeSale();
        timeSale.setId(timeSaleId);
        timeSale.setQuantity(10L);
        timeSale.setRemainingQuantity(10L);
        timeSale.setDiscountPrice(800L);
        timeSale.setStartAt(LocalDateTime.now().minusMinutes(10));
        timeSale.setEndAt(LocalDateTime.now().plusHours(1));
        timeSale.setStatus(TimeSale.Status.ACTIVE);

        when(redissonClient.getLock("time-sale-lock:" + timeSaleId)).thenReturn(rLock);
        when(rLock.tryLock(3L, 3L, TimeUnit.SECONDS)).thenReturn(true);
        when(redissonClient.getBucket("time-sale:" + timeSaleId)).thenReturn(rBucket);
        when(rBucket.get()).thenReturn(null);
        when(timeSaleRepository.findById(timeSaleId)).thenReturn(Optional.of(timeSale));

        IllegalStateException ex = assertThrows(IllegalStateException.class,
                () -> timeSaleService.purchase(request, userId));
        assertEquals("Not enough quantity available", ex.getMessage());
        verify(rLock).unlock();
    }
}
