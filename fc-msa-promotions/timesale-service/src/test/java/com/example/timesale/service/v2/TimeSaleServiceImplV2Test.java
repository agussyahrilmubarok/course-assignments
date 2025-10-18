package com.example.timesale.service.v2;

import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.redisson.api.RBucket;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;

import java.time.LocalDateTime;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class TimeSaleServiceImplV2Test {

    private static final String TEST_TS_ID = "TS-V2-001";
    private static final String USER_ID = "USER_A";
    @InjectMocks
    private TimeSaleServiceImplV2 timeSaleService;
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
    private RBucket<Object> mockBucket;
    @Mock
    private RLock mockLock;
    private TimeSaleDTO.CreateRequest createRequest;
    private TimeSaleDTO.PurchaseRequest purchaseRequest;
    private Product product;
    private TimeSale timeSale;

    @BeforeEach
    void setUp() throws Exception {
        createRequest = TimeSaleDTO.CreateRequest.builder()
                .quantity(5L)
                .discountPrice(2000L)
                .startAt(LocalDateTime.now().minusHours(1))
                .endAt(LocalDateTime.now().plusHours(2))
                .product(TimeSaleDTO.ProductRequest.builder()
                        .name("Product V2")
                        .price(5000L)
                        .build())
                .build();

        purchaseRequest = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(TEST_TS_ID)
                .quantity(2L)
                .build();

        product = new Product();
        product.setId(UUID.randomUUID().toString());
        product.setName("Product V2");
        product.setPrice(5000L);

        timeSale = new TimeSale();
        timeSale.setId(TEST_TS_ID);
        timeSale.setQuantity(5L);
        timeSale.setRemainingQuantity(5L);
        timeSale.setDiscountPrice(2000L);
        timeSale.setStartAt(createRequest.getStartAt());
        timeSale.setEndAt(createRequest.getEndAt());
        timeSale.setStatus(TimeSale.Status.ACTIVE);
        timeSale.setProduct(product);

        when(redissonClient.getBucket("time-sale:" + TEST_TS_ID)).thenReturn(mockBucket);
    }

    @Test
    void shouldCreateTimeSaleAndCacheIt() throws Exception {
        String json = "{\"id\":\"" + TEST_TS_ID + "\",\"quantity\":5}";
        when(objectMapper.writeValueAsString(any(TimeSale.class))).thenReturn(json);
        when(productRepository.save(any())).thenReturn(product);
        when(timeSaleRepository.save(any())).thenAnswer(inv -> {
            TimeSale ts = inv.getArgument(0);
            ts.setId(TEST_TS_ID);
            return ts;
        });

        TimeSale result = timeSaleService.create(createRequest);

        assertNotNull(result);
        assertEquals(TEST_TS_ID, result.getId());
        verify(mockBucket).set(eq(json));
    }

    @Test
    void shouldFindByIdFromCache() throws Exception {
        String json = "{\"id\":\"" + TEST_TS_ID + "\"}";
        when(mockBucket.get()).thenReturn(json);
        when(objectMapper.readValue(eq(json), eq(TimeSale.class))).thenReturn(timeSale);

        TimeSale result = timeSaleService.findById(TEST_TS_ID);

        assertNotNull(result);
        assertEquals(TEST_TS_ID, result.getId());
    }

    @Test
    void shouldFallbackToRepositoryWhenCacheMiss() throws Exception {
        when(mockBucket.get()).thenReturn(null);
        when(timeSaleRepository.findById(TEST_TS_ID)).thenReturn(Optional.of(timeSale));
        String json = "{\"id\":\"" + TEST_TS_ID + "\"}";
        when(objectMapper.writeValueAsString(timeSale)).thenReturn(json);

        TimeSale result = timeSaleService.findById(TEST_TS_ID);

        assertNotNull(result);
        verify(mockBucket).set(eq(json));
    }

    @Test
    void shouldThrowIfNotFoundInRepository() {
        when(mockBucket.get()).thenReturn(null);
        when(timeSaleRepository.findById(TEST_TS_ID)).thenReturn(Optional.empty());

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class,
                () -> timeSaleService.findById(TEST_TS_ID));
        assertEquals("TimeSale not found", ex.getMessage());
    }

    @Test
    void shouldPurchaseWithLock() throws Exception {
        when(redissonClient.getLock("time-sale-lock:" + TEST_TS_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockBucket.get()).thenReturn("{\"id\":\"" + TEST_TS_ID + "\"}");
        when(objectMapper.readValue(anyString(), eq(TimeSale.class))).thenReturn(timeSale);
        when(timeSaleOrderRepository.save(any())).thenAnswer(inv -> inv.getArgument(0));
        when(timeSaleRepository.save(any())).thenReturn(timeSale);

        TimeSale result = timeSaleService.purchase(purchaseRequest, USER_ID);

        assertNotNull(result);
        assertEquals(3L, result.getRemainingQuantity());
        verify(mockBucket, atLeastOnce()).set(any());
        verify(mockLock).unlock();
    }

//    @Test
//    void shouldThrowWhenCannotAcquireLock() throws Exception {
//        when(redissonClient.getLock("time-sale-lock:" + TEST_TS_ID)).thenReturn(mockLock);
//        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(false);
//
//        RuntimeException ex = assertThrows(RuntimeException.class,
//                () -> timeSaleService.purchase(purchaseRequest, USER_ID));
//        assertEquals("Could not acquire lock for time sale purchase.", ex.getMessage());
//    }

    @Test
    void shouldThrowOnUnexpectedErrorAndReleaseLock() throws Exception {
        when(redissonClient.getLock("time-sale-lock:" + TEST_TS_ID)).thenReturn(mockLock);
        when(mockLock.tryLock(anyLong(), anyLong(), any())).thenReturn(true);
        when(mockBucket.get()).thenReturn("{\"id\":\"" + TEST_TS_ID + "\"}");
        when(objectMapper.readValue(anyString(), eq(TimeSale.class))).thenReturn(timeSale);
        when(timeSaleRepository.save(any())).thenThrow(new RuntimeException("db error"));

        RuntimeException ex = assertThrows(RuntimeException.class,
                () -> timeSaleService.purchase(purchaseRequest, USER_ID));
        assertEquals("Failed to process purchase", ex.getMessage());
        verify(mockLock).unlock();
    }
}
