package com.finly.finly.controller;

import com.finly.finly.model.SignUpDTO;
import com.finly.finly.model.UserDTO;
import com.finly.finly.service.UserService;
import com.finly.finly.util.WebUtils;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.util.Map;
import java.util.UUID;

@Controller
@Slf4j
@RequestMapping("/signup")
@RequiredArgsConstructor
public class SignUpController extends BaseController {

    private final UserService userService;

    @GetMapping
    public String signUp(Model model) {
        if (isLoggedIn()) {
            log.debug("Authenticated user attempted to access signup page. Redirecting to dashboard.");
            return "redirect:/dashboard";
        }

        model.addAttribute("dto", new SignUpDTO());
        log.debug("Displaying signup page.");
        return "home/signup";
    }

    @PostMapping
    public String signUp(@ModelAttribute("dto") @Valid SignUpDTO signUpDTO,
                         BindingResult bindingResult,
                         RedirectAttributes redirectAttributes) {
        log.info("Processing signup attempt for email: {}", signUpDTO.getEmail());

        if (bindingResult.hasErrors()) {
            log.warn("Signup validation failed for email {}: {}", signUpDTO.getEmail(), bindingResult.getAllErrors());
            return "home/signup";
        }

        try {
            UserDTO param = new UserDTO();
            param.setEmail(signUpDTO.getEmail());
            param.setPassword(signUpDTO.getPassword());
            UUID createdId = userService.create(param);

            log.info("Signup successful. New user ID: {}", createdId);
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("home.signup.success"),
                    "type", WebUtils.MSG_SUCCESS
            ));
            return "redirect:/login";

        } catch (Exception ex) {
            log.error("Signup failed for email {}: {}", signUpDTO.getEmail(), ex.getMessage());
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("home.signup.failed"),
                    "type", WebUtils.MSG_ERROR
            ));
            return "redirect:/signup";
        }
    }
}
