package com.finly.finly.controller;

import com.finly.finly.model.InvoiceDTO;
import com.finly.finly.service.CustomerService;
import com.finly.finly.service.InvoiceService;
import com.finly.finly.util.WebUtils;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.util.Map;
import java.util.UUID;


@Controller
@RequestMapping("/dashboard/invoices")
@Slf4j
@RequiredArgsConstructor
public class InvoiceController extends BaseController {

    private final InvoiceService invoiceService;
    private final CustomerService customerService;

    @GetMapping
    public String index(@RequestParam(value = "search", required = false) String search,
                        final Model model) {
        log.info("Retrieving invoices for user ID: {}. Search term: {}", getUserId(), search != null ? search : "none");

        model.addAttribute("typePage", "data");
        model.addAttribute("invoices", invoiceService.findInvoiceByOwnerAndSearch(getUserId(), search));

        log.debug("Invoices data prepared for rendering.");
        return "invoice/index";
    }

    @GetMapping("/create")
    public String create(final Model model) {
        log.debug("Rendering create invoice form for user ID: {}", getUserId());

        model.addAttribute("typePage", "form");
        model.addAttribute("customers", customerService.findByUserId(getUserId()));
        model.addAttribute("invoice", new InvoiceDTO());

        log.debug("Customers list and invoice form initialized for creation.");
        return "invoice/index";
    }

    @PostMapping
    public String save(@ModelAttribute("invoice") @Valid InvoiceDTO invoiceDTO,
                       BindingResult bindingResult,
                       RedirectAttributes redirectAttributes,
                       final Model model) {
        log.info("Processing creation of new invoice for user ID: {}", getUserId());
        if (bindingResult.hasErrors()) {
            log.warn("Invoice creation validation failed for user ID {}: {}", getUserId(), bindingResult.getAllErrors());
            model.addAttribute("typePage", "form");
            model.addAttribute("customers", customerService.findByUserId(getUserId()));
            model.addAttribute("invoice", invoiceDTO);
            return "invoice/index";
        }

        invoiceDTO.setOwner(getUserId());
        UUID createdId = invoiceService.create(invoiceDTO);
        redirectAttributes.addFlashAttribute("info", Map.of(
                "message", WebUtils.getMessage("invoice.create.success"),
                "type", WebUtils.MSG_SUCCESS
        ));

        log.info("Invoice created successfully with ID: {}", createdId);
        return "redirect:/dashboard/invoices";
    }

    @GetMapping("/{id}/edit")
    public String edit(@PathVariable("id") UUID id, final Model model) {
        log.debug("Rendering edit form for invoice ID: {}", id);

        model.addAttribute("title", "Edit Invoice");
        model.addAttribute("typePage", "form");
        model.addAttribute("customers", customerService.findByUserId(getUserId()));
        model.addAttribute("invoice", invoiceService.getDetail(id));

        log.debug("Invoice data loaded for invoice ID: {}", id);
        return "invoice/index";
    }

    @PostMapping("/{id}")
    public String update(@PathVariable("id") UUID id,
                         @ModelAttribute("invoice") @Valid InvoiceDTO invoiceDTO,
                         BindingResult bindingResult,
                         RedirectAttributes redirectAttributes,
                         final Model model) {
        log.info("Processing update for invoice ID: {}", id);

        if (bindingResult.hasErrors()) {
            log.warn("Invoice update validation failed for invoice ID {}: {}", id, bindingResult.getAllErrors());
            model.addAttribute("title", "Edit Invoice");
            model.addAttribute("typePage", "form");
            model.addAttribute("invoice", invoiceDTO);
            return "invoice/index";
        }

        invoiceDTO.setOwner(getUserId());
        invoiceService.update(id, invoiceDTO);
        redirectAttributes.addFlashAttribute("info", Map.of(
                "message", WebUtils.getMessage("invoice.update.success"),
                "type", WebUtils.MSG_SUCCESS
        ));

        log.info("Invoice updated successfully with ID: {}", id);
        return "redirect:/dashboard/invoices";
    }

    @PostMapping("/{id}/delete")
    public String delete(@PathVariable("id") UUID id, RedirectAttributes redirectAttributes) {
        log.info("Processing deletion of invoice ID: {}", id);
        try {
            invoiceService.delete(id);
            log.info("Invoice deleted successfully with ID: {}", id);
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("invoice.delete.success"),
                    "type", WebUtils.MSG_SUCCESS
            ));
        } catch (Exception e) {
            log.error("Failed to delete invoice with ID: {}", id, e);
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("invoice.delete.failed"),
                    "type", WebUtils.MSG_ERROR
            ));
        }
        return "redirect:/dashboard/invoices";
    }
}