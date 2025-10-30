package com.example.timesale.service.v1;

import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.domain.TimeSaleOrder;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class TimeSaleServiceImplV1Test {

    @InjectMocks
    private TimeSaleServiceImplV1 timeSaleServiceImplV1;

    @Mock
    private TimeSaleRepository timeSaleRepository;
    @Mock
    private ProductRepository productRepository;
    @Mock
    private TimeSaleOrderRepository timeSaleOrderRepository;

    @Test
    void testCreate_whenValidRequest_shouldReturnSavedTimeSale() {
        TimeSaleDTO.ProductRequest productRequest = TimeSaleDTO.ProductRequest.builder()
                .name("Test Product")
                .price(1000L)
                .build();

        TimeSaleDTO.CreateRequest request = TimeSaleDTO.CreateRequest.builder()
                .product(productRequest)
                .quantity(10L)
                .discountPrice(800L)
                .startAt(LocalDateTime.now().minusMinutes(1))
                .endAt(LocalDateTime.now().plusHours(1))
                .build();

        Product savedProduct = new Product();
        savedProduct.setId("prod-1");
        savedProduct.setName("Test Product");
        savedProduct.setPrice(1000L);

        TimeSale savedTimeSale = new TimeSale();
        savedTimeSale.setId("ts-1");
        savedTimeSale.setQuantity(10L);
        savedTimeSale.setRemainingQuantity(10L);
        savedTimeSale.setDiscountPrice(800L);
        savedTimeSale.setStartAt(request.getStartAt());
        savedTimeSale.setEndAt(request.getEndAt());
        savedTimeSale.setProduct(savedProduct);
        savedTimeSale.setStatus(TimeSale.Status.ACTIVE);

        when(productRepository.save(any(Product.class))).thenReturn(savedProduct);
        when(timeSaleRepository.save(any(TimeSale.class))).thenReturn(savedTimeSale);

        TimeSale result = timeSaleServiceImplV1.create(request);

        assertNotNull(result);
        assertEquals(10L, result.getQuantity());
        verify(productRepository).save(any(Product.class));
        verify(timeSaleRepository).save(any(TimeSale.class));
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

        Exception exception = assertThrows(IllegalArgumentException.class, () ->
                timeSaleServiceImplV1.create(request)
        );

        assertEquals("Start time must be before end time", exception.getMessage());
    }

    @Test
    void testFindById_whenExists_shouldReturnTimeSale() {
        TimeSale timeSale = new TimeSale();
        timeSale.setId("ts-123");

        when(timeSaleRepository.findById("ts-123")).thenReturn(Optional.of(timeSale));

        TimeSale result = timeSaleServiceImplV1.findById("ts-123");

        assertNotNull(result);
        assertEquals("ts-123", result.getId());
    }

    @Test
    void testFindById_whenNotFound_shouldThrowException() {
        when(timeSaleRepository.findById("invalid-id")).thenReturn(Optional.empty());

        Exception exception = assertThrows(IllegalArgumentException.class, () ->
                timeSaleServiceImplV1.findById("invalid-id")
        );

        assertEquals("Time sale not found", exception.getMessage());
    }

    @Test
    void testFindAllOngoing_whenCalled_shouldReturnPageOfTimeSales() {
        LocalDateTime now = LocalDateTime.now();
        Pageable pageable = PageRequest.of(0, 10);

        TimeSale ts = new TimeSale();
        ts.setId("ts-1");

        Page<TimeSale> expectedPage = new PageImpl<>(List.of(ts));

        when(timeSaleRepository.findAllByStartAtBeforeAndEndAtAfterAndStatus(
                any(LocalDateTime.class), eq(TimeSale.Status.ACTIVE), eq(pageable)
        )).thenReturn(expectedPage);

        Page<TimeSale> result = timeSaleServiceImplV1.findAllOngoing(pageable);

        assertNotNull(result);
        assertEquals(1, result.getContent().size());
    }

    @Test
    void testPurchase_whenValidRequest_shouldReturnUpdatedTimeSale() {
        String userId = "user-1";
        String timeSaleId = "ts-1";

        TimeSaleDTO.PurchaseRequest request = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(timeSaleId)
                .quantity(2L)
                .build();

        Product product = new Product();
        product.setId("prod-1");
        product.setName("Product 1");
        product.setPrice(1000L);

        TimeSale timeSale = new TimeSale();
        timeSale.setId(timeSaleId);
        timeSale.setQuantity(10L);
        timeSale.setRemainingQuantity(10L);
        timeSale.setDiscountPrice(800L);
        timeSale.setStartAt(LocalDateTime.now().minusMinutes(10));
        timeSale.setEndAt(LocalDateTime.now().plusHours(1));
        timeSale.setStatus(TimeSale.Status.ACTIVE);
        timeSale.setProduct(product);

        when(timeSaleRepository.findByIdWithPessimisticLock(timeSaleId)).thenReturn(Optional.of(timeSale));
        when(timeSaleRepository.save(any(TimeSale.class))).thenAnswer(invocation -> invocation.getArgument(0));
        when(timeSaleOrderRepository.save(any(TimeSaleOrder.class))).thenReturn(new TimeSaleOrder());

        TimeSale result = timeSaleServiceImplV1.purchase(request, userId);

        assertNotNull(result);
        assertEquals(8L, result.getRemainingQuantity());
        verify(timeSaleOrderRepository).save(any(TimeSaleOrder.class));
    }

    @Test
    void testPurchase_whenTimeSaleNotFound_shouldThrowException() {
        String timeSaleId = "invalid-id";
        TimeSaleDTO.PurchaseRequest request = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId(timeSaleId)
                .quantity(1L)
                .build();

        when(timeSaleRepository.findByIdWithPessimisticLock(timeSaleId)).thenReturn(Optional.empty());

        Exception exception = assertThrows(IllegalArgumentException.class, () ->
                timeSaleServiceImplV1.purchase(request, "user-123")
        );

        assertEquals("TimeSale not found", exception.getMessage());
    }

    @Test
    void testCreate_whenQuantityIsZero_shouldThrowException() {
        TimeSaleDTO.CreateRequest request = TimeSaleDTO.CreateRequest.builder()
                .product(TimeSaleDTO.ProductRequest.builder().name("Product A").price(1000L).build())
                .quantity(0L)  // invalid
                .discountPrice(900L)
                .startAt(LocalDateTime.now().minusMinutes(1))
                .endAt(LocalDateTime.now().plusHours(1))
                .build();

        Exception exception = assertThrows(IllegalArgumentException.class, () ->
                timeSaleServiceImplV1.create(request)
        );

        assertEquals("Quantity must be positive", exception.getMessage());
    }

    @Test
    void testCreate_whenDiscountPriceIsZero_shouldThrowException() {
        TimeSaleDTO.CreateRequest request = TimeSaleDTO.CreateRequest.builder()
                .product(TimeSaleDTO.ProductRequest.builder().name("Product A").price(1000L).build())
                .quantity(5L)
                .discountPrice(0L)  // invalid
                .startAt(LocalDateTime.now().minusMinutes(1))
                .endAt(LocalDateTime.now().plusHours(1))
                .build();

        Exception exception = assertThrows(IllegalArgumentException.class, () ->
                timeSaleServiceImplV1.create(request)
        );

        assertEquals("Discount price must be positive", exception.getMessage());
    }
}
