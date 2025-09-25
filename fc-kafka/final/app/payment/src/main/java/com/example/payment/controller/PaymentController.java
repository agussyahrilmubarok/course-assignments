package com.example.payment.controller;

import com.example.payment.model.PaymentDTO;
import com.example.payment.service.PaymentService;
import com.example.payment.util.WebUtils;
import jakarta.validation.Valid;
import java.util.UUID;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;


@Controller
@RequestMapping("/payments")
public class PaymentController {

    private final PaymentService paymentService;

    public PaymentController(final PaymentService paymentService) {
        this.paymentService = paymentService;
    }

    @GetMapping
    public String list(final Model model) {
        model.addAttribute("payments", paymentService.findAll());
        return "payment/list";
    }

    @GetMapping("/add")
    public String add(@ModelAttribute("payment") final PaymentDTO paymentDTO) {
        return "payment/add";
    }

    @PostMapping("/add")
    public String add(@ModelAttribute("payment") @Valid final PaymentDTO paymentDTO,
            final BindingResult bindingResult, final RedirectAttributes redirectAttributes) {
        if (bindingResult.hasErrors()) {
            return "payment/add";
        }
        paymentService.create(paymentDTO);
        redirectAttributes.addFlashAttribute(WebUtils.MSG_SUCCESS, WebUtils.getMessage("payment.create.success"));
        return "redirect:/payments";
    }

    @GetMapping("/edit/{id}")
    public String edit(@PathVariable(name = "id") final UUID id, final Model model) {
        model.addAttribute("payment", paymentService.get(id));
        return "payment/edit";
    }

    @PostMapping("/edit/{id}")
    public String edit(@PathVariable(name = "id") final UUID id,
            @ModelAttribute("payment") @Valid final PaymentDTO paymentDTO,
            final BindingResult bindingResult, final RedirectAttributes redirectAttributes) {
        if (bindingResult.hasErrors()) {
            return "payment/edit";
        }
        paymentService.update(id, paymentDTO);
        redirectAttributes.addFlashAttribute(WebUtils.MSG_SUCCESS, WebUtils.getMessage("payment.update.success"));
        return "redirect:/payments";
    }

    @PostMapping("/delete/{id}")
    public String delete(@PathVariable(name = "id") final UUID id,
            final RedirectAttributes redirectAttributes) {
        paymentService.delete(id);
        redirectAttributes.addFlashAttribute(WebUtils.MSG_INFO, WebUtils.getMessage("payment.delete.success"));
        return "redirect:/payments";
    }

}
