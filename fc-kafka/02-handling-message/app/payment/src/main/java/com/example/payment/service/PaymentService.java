package com.example.payment.service;

import com.example.payment.domain.Payment;
import com.example.payment.model.PaymentDTO;
import com.example.payment.repos.PaymentRepository;
import com.example.payment.util.NotFoundException;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.time.ZoneOffset;
import java.util.List;
import java.util.UUID;


@Service
public class PaymentService {

    private static final ZoneOffset DEFAULT_ZONE_OFFSET = ZoneOffset.ofHours(7); // Asia/Jakarta

    private final PaymentRepository paymentRepository;

    public PaymentService(final PaymentRepository paymentRepository) {
        this.paymentRepository = paymentRepository;
    }

    public List<PaymentDTO> findAll() {
        final List<Payment> payments = paymentRepository.findAll(Sort.by("id"));
        return payments.stream()
                .map(payment -> mapToDTO(payment, new PaymentDTO()))
                .toList();
    }

    public PaymentDTO get(final UUID id) {
        return paymentRepository.findById(id)
                .map(payment -> mapToDTO(payment, new PaymentDTO()))
                .orElseThrow(NotFoundException::new);
    }

    public UUID create(final PaymentDTO paymentDTO) {
        final Payment payment = new Payment();
        mapToEntity(paymentDTO, payment);
        return paymentRepository.save(payment).getId();
    }

    public void update(final UUID id, final PaymentDTO paymentDTO) {
        final Payment payment = paymentRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        mapToEntity(paymentDTO, payment);
        paymentRepository.save(payment);
    }

    public void delete(final UUID id) {
        final Payment payment = paymentRepository.findById(id)
                .orElseThrow(NotFoundException::new);
        paymentRepository.delete(payment);
    }

    private PaymentDTO mapToDTO(final Payment payment, final PaymentDTO paymentDTO) {
        paymentDTO.setId(payment.getId());
        paymentDTO.setOrderId(payment.getOrderId());
        paymentDTO.setCustomerId(payment.getCustomerId());
        paymentDTO.setAmount(payment.getAmount());
        paymentDTO.setMethod(payment.getMethod());
        paymentDTO.setStatus(payment.getStatus() != null ? payment.getStatus().name() : null);
        paymentDTO.setPaidAt(payment.getPaidAt() != null ? payment.getPaidAt().toLocalDateTime() : null);
        return paymentDTO;
    }

    private Payment mapToEntity(final PaymentDTO paymentDTO, final Payment payment) {
        payment.setOrderId(paymentDTO.getOrderId());
        payment.setCustomerId(paymentDTO.getCustomerId());
        payment.setAmount(paymentDTO.getAmount());
        payment.setMethod(paymentDTO.getMethod());
        if (paymentDTO.getStatus() != null) {
            payment.setStatus(Payment.Status.valueOf(paymentDTO.getStatus()));
        }
        payment.setPaidAt(paymentDTO.getPaidAt() != null ? paymentDTO.getPaidAt().atOffset(DEFAULT_ZONE_OFFSET) : null);
        return payment;
    }

}
