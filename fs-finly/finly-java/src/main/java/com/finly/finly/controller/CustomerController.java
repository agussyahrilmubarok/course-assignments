package com.finly.finly.controller;

import com.finly.finly.model.CustomerDTO;
import com.finly.finly.service.CustomerService;
import com.finly.finly.util.WebUtils;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.util.List;
import java.util.Map;
import java.util.UUID;

@Controller
@RequestMapping("/dashboard/customers")
@Slf4j
@RequiredArgsConstructor
public class CustomerController extends BaseController {

    private final CustomerService customerService;

    @GetMapping
    public String index(@RequestParam(value = "search", required = false) String search,
                        final Model model) {
        log.info("Listing customers. Search query: {}", search != null ? search : "none");

        List<CustomerDTO> customers = customerService.findCustomersByOwnerAndSearch(getUserId(), search);
        model.addAttribute("typePage", "data");
        model.addAttribute("customers", customers);

        return "customer/index";
    }

    @GetMapping("/create")
    public String create(final Model model) {
        log.debug("Displaying form to create a new customer.");

        model.addAttribute("typePage", "form");
        model.addAttribute("customer", new CustomerDTO());

        return "customer/index";
    }

    @PostMapping
    public String save(@ModelAttribute("customer") @Valid CustomerDTO customerDTO,
                       BindingResult bindingResult,
                       RedirectAttributes redirectAttributes,
                       final Model model) {
        if (bindingResult.hasErrors()) {
            log.warn("Customer creation validation failed: {}", bindingResult.getAllErrors());
            model.addAttribute("typePage", "form");
            model.addAttribute("customer", customerDTO);
            return "customer/index";
        }

        customerDTO.setUser(getUserId());
        UUID createdId = customerService.create(customerDTO);

        log.info("Customer created successfully. ID: {}", createdId);
        redirectAttributes.addFlashAttribute("info", Map.of(
                "message", WebUtils.getMessage("customer.create.success"),
                "type", WebUtils.MSG_SUCCESS
        ));

        return "redirect:/dashboard/customers";
    }

    @GetMapping("/{id}/edit")
    public String edit(@PathVariable("id") UUID id,
                       final Model model) {
        log.debug("Displaying form to edit customer with ID: {}", id);

        model.addAttribute("typePage", "form");
        model.addAttribute("customer", customerService.get(id));

        return "customer/index";
    }

    @PostMapping("/{id}")
    public String update(@PathVariable("id") UUID id,
                         @ModelAttribute("customer") @Valid CustomerDTO customerDTO,
                         BindingResult bindingResult,
                         RedirectAttributes redirectAttributes,
                         final Model model) {
        if (bindingResult.hasErrors()) {
            log.warn("Customer update validation failed for ID {}: {}", id, bindingResult.getAllErrors());
            model.addAttribute("typePage", "form");
            model.addAttribute("customer", customerDTO);
            return "customer/index";
        }

        customerDTO.setUser(getUserId());
        customerService.update(id, customerDTO);
        log.info("Customer with ID {} updated successfully.", id);

        redirectAttributes.addFlashAttribute("info", Map.of(
                "message", WebUtils.getMessage("customer.update.success"),
                "type", WebUtils.MSG_SUCCESS
        ));

        return "redirect:/dashboard/customers";
    }

    @PostMapping("/{id}/delete")
    public String delete(@PathVariable("id") UUID id,
                         RedirectAttributes redirectAttributes) {
        try {
            customerService.delete(id);
            log.info("Customer with ID {} deleted successfully.", id);
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("customer.delete.success"),
                    "type", WebUtils.MSG_SUCCESS
            ));
        } catch (Exception e) {
            log.error("Failed to delete customer with ID {}.", id, e);
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("customer.delete.failed"),
                    "type", WebUtils.MSG_ERROR
            ));
        }

        return "redirect:/dashboard/customers";
    }
}