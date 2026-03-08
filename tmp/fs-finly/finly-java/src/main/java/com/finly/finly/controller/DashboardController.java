package com.finly.finly.controller;

import com.finly.finly.model.InvoiceDetailDTO;
import com.finly.finly.service.CustomerService;
import com.finly.finly.service.InvoiceService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;

import java.text.NumberFormat;
import java.time.LocalDate;
import java.util.*;

@Controller
@Slf4j
@RequiredArgsConstructor
public class DashboardController extends BaseController {

    private final CustomerService customerService;
    private final InvoiceService invoiceService;

    @GetMapping("/dashboard")
    public String index(final Model model) {
        if (!isLoggedIn()) return "redirect:/login";

        var userId = getUserId();

        // Retrieve counts
        var customerCount = customerService.countByUser(userId);
        var invoiceCount = invoiceService.countByOwner(userId);
        log.debug("Dashboard counts retrieved for user.", Map.of(
                "userId", userId,
                "customerCount", customerCount,
                "invoiceCount", invoiceCount
        ));

        // Retrieve all invoices
        List<InvoiceDetailDTO> allInvoices = invoiceService.findInvoiceByOwnerAndSearch(userId, "");
        log.debug("All invoices retrieved for dashboard.", Map.of(
                "userId", userId,
                "totalInvoices", allInvoices.size()
        ));

        // Calculate totals
        double totalPaid = allInvoices.stream()
                .filter(invoice -> "paid".equalsIgnoreCase(invoice.getStatus()))
                .mapToDouble(i -> i.getAmount().doubleValue())
                .sum();
        double totalPending = allInvoices.stream()
                .filter(invoice -> "pending".equalsIgnoreCase(invoice.getStatus()))
                .mapToDouble(i -> i.getAmount().doubleValue())
                .sum();
        log.debug("Dashboard totals calculated.", Map.of(
                "userId", userId,
                "totalPaid", totalPaid,
                "totalPending", totalPending
        ));

        // Latest invoices
        List<InvoiceDetailDTO> latestInvoices = allInvoices.stream()
                .sorted(Comparator.comparing(InvoiceDetailDTO::getDueDate).reversed())
                .limit(5)
                .toList();
        log.debug("Latest 5 invoices retrieved.", Map.of(
                "userId", userId,
                "latestInvoicesCount", latestInvoices.size()
        ));

        // Revenue for last 6 months
        List<Map<String, Object>> revenueData = new ArrayList<>();
        for (int i = 0; i < 6; i++) {
            LocalDate monthDate = LocalDate.now().minusMonths(i);
            String monthLabel = monthDate.getMonth().toString().substring(0, 3);
            double revenueForMonth = allInvoices.stream()
                    .filter(invoice -> invoice.getDueDate().getMonth() == monthDate.getMonth())
                    .mapToDouble(a -> a.getAmount().doubleValue())
                    .sum();

            Map<String, Object> revenueEntry = new HashMap<>();
            revenueEntry.put("month", monthLabel);
            revenueEntry.put("revenue", revenueForMonth);
            revenueData.add(0, revenueEntry);
        }
        log.info("Revenue data for last 6 months prepared for user ID {}.", userId);

        // Set model attributes
        model.addAttribute("latestInvoices", latestInvoices);
        model.addAttribute("revenueData", revenueData);
        model.addAttribute("invoiceCount", invoiceCount);
        model.addAttribute("customerCount", customerCount);
        model.addAttribute("totalPaid", totalPaid);
        model.addAttribute("totalPending", totalPending);

        return "dashboard/index";
    }

    private LocalDate convertToLocalDate(Date date) {
        return date.toInstant()
                .atZone(java.time.ZoneId.systemDefault())
                .toLocalDate();
    }

    private String formatUSD(double amount) {
        NumberFormat usdFormatter = NumberFormat.getCurrencyInstance(Locale.US);
        return usdFormatter.format(amount);
    }
}