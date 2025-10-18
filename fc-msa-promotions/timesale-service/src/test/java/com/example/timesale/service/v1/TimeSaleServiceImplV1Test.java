package com.example.timesale.service.v1;

import com.example.timesale.domain.Product;
import com.example.timesale.domain.TimeSale;
import com.example.timesale.model.TimeSaleDTO;
import com.example.timesale.repos.ProductRepository;
import com.example.timesale.repos.TimeSaleOrderRepository;
import com.example.timesale.repos.TimeSaleRepository;
import org.junit.jupiter.api.BeforeEach;
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
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class TimeSaleServiceImplV1Test {

    @InjectMocks
    private TimeSaleServiceImplV1 timeSaleService;

    @Mock
    private TimeSaleRepository timeSaleRepository;

    @Mock
    private ProductRepository productRepository;

    @Mock
    private TimeSaleOrderRepository timeSaleOrderRepository;

    private TimeSaleDTO.CreateRequest createRequest;
    private TimeSaleDTO.PurchaseRequest purchaseRequest;
    private Product product;
    private TimeSale timeSale;

    @BeforeEach
    void setUp() {
        createRequest = TimeSaleDTO.CreateRequest.builder()
                .quantity(10L)
                .discountPrice(5000L)
                .startAt(LocalDateTime.now().minusHours(1))
                .endAt(LocalDateTime.now().plusHours(1))
                .product(TimeSaleDTO.ProductRequest.builder()
                        .name("Test Product")
                        .price(10000L)
                        .build())
                .build();

        purchaseRequest = TimeSaleDTO.PurchaseRequest.builder()
                .timeSaleId("TS001")
                .quantity(2L)
                .build();

        product = new Product();
        product.setId(UUID.randomUUID().toString());
        product.setName("Test Product");
        product.setPrice(10000L);

        timeSale = new TimeSale();
        timeSale.setId("TS001");
        timeSale.setQuantity(10L);
        timeSale.setRemainingQuantity(10L);
        timeSale.setDiscountPrice(5000L);
        timeSale.setStartAt(LocalDateTime.now().minusHours(1));
        timeSale.setEndAt(LocalDateTime.now().plusHours(1));
        timeSale.setStatus(TimeSale.Status.ACTIVE);
        timeSale.setProduct(product);
    }

    @Test
    void shouldCreateTimeSaleSuccessfully() {
        when(productRepository.save(any())).thenReturn(product);
        when(timeSaleRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

        TimeSale result = timeSaleService.create(createRequest);

        assertNotNull(result);
        assertThat(result.getQuantity()).isEqualTo(10L);
        assertThat(result.getDiscountPrice()).isEqualTo(5000L);
        verify(productRepository).save(any());
        verify(timeSaleRepository).save(any());
    }

    @Test
    void shouldThrowExceptionWhenStartAfterEnd() {
        createRequest.setStartAt(LocalDateTime.now().plusDays(1));
        createRequest.setEndAt(LocalDateTime.now());

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> timeSaleService.create(createRequest));
        assertThat(ex.getMessage()).isEqualTo("Start time must be before end time");
    }

    @Test
    void shouldThrowExceptionWhenQuantityInvalid() {
        createRequest.setQuantity(0L);

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> timeSaleService.create(createRequest));
        assertThat(ex.getMessage()).isEqualTo("Quantity must be positive");
    }

    @Test
    void shouldThrowExceptionWhenDiscountPriceInvalid() {
        createRequest.setDiscountPrice(0L);

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> timeSaleService.create(createRequest));
        assertThat(ex.getMessage()).isEqualTo("Discount price must be positive");
    }

    @Test
    void shouldFindByIdSuccessfully() {
        when(timeSaleRepository.findById("TS001")).thenReturn(Optional.of(timeSale));

        TimeSale result = timeSaleService.findById("TS001");

        assertNotNull(result);
        assertEquals("TS001", result.getId());
    }

    @Test
    void shouldThrowExceptionWhenFindByIdNotFound() {
        when(timeSaleRepository.findById("TS001")).thenReturn(Optional.empty());

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () -> timeSaleService.findById("TS001"));
        assertThat(ex.getMessage()).isEqualTo("Time sale not found");
    }

    @Test
    void shouldFindAllOngoingSales() {
        Pageable pageable = PageRequest.of(0, 10);
        Page<TimeSale> page = new PageImpl<>(List.of(timeSale));

        when(timeSaleRepository.findAllByStartAtBeforeAndEndAtAfterAndStatus(
                any(), eq(TimeSale.Status.ACTIVE), eq(pageable)))
                .thenReturn(page);

        Page<TimeSale> result = timeSaleService.findAllOngoing(pageable);

        assertThat(result.getContent()).hasSize(1);
        assertThat(result.getContent().get(0).getId()).isEqualTo("TS001");
    }

    @Test
    void shouldPurchaseSuccessfully() {
        when(timeSaleRepository.findByIdWithPessimisticLock("TS001")).thenReturn(Optional.of(timeSale));
        when(timeSaleRepository.save(any())).thenReturn(timeSale);
        when(timeSaleOrderRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

        TimeSale result = timeSaleService.purchase(purchaseRequest, "USER123");

        assertNotNull(result);
        assertThat(result.getRemainingQuantity()).isLessThanOrEqualTo(10L);
        verify(timeSaleRepository).save(any());
        verify(timeSaleOrderRepository).save(any());
    }

    @Test
    void shouldThrowExceptionWhenPurchaseNotFound() {
        when(timeSaleRepository.findByIdWithPessimisticLock("TS001")).thenReturn(Optional.empty());

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, () ->
                timeSaleService.purchase(purchaseRequest, "USER123"));

        assertThat(ex.getMessage()).isEqualTo("TimeSale not found");
    }
}
