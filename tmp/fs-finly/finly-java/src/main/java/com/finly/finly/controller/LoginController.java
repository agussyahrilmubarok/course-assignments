package com.finly.finly.controller;

import com.finly.finly.model.LoginDTO;
import com.finly.finly.util.WebUtils;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.util.Map;

@Controller
@Slf4j
@RequestMapping("/login")
@RequiredArgsConstructor
public class LoginController extends BaseController {

    private final AuthenticationManager authenticationManager;

    @GetMapping
    public String login(Model model) {
        if (isLoggedIn()) {
            log.debug("User is already authenticated. Redirecting to dashboard.");
            return "redirect:/dashboard";
        }

        model.addAttribute("dto", new LoginDTO());
        log.debug("Displaying login page.");
        return "home/login";
    }

    @PostMapping
    public String login(@ModelAttribute("dto") @Valid LoginDTO loginDTO,
                        BindingResult bindingResult,
                        RedirectAttributes redirectAttributes) {
        log.info("Processing login attempt for email: {}", loginDTO.getEmail());

        if (bindingResult.hasErrors()) {
            log.warn("Login validation failed for email {}: {}", loginDTO.getEmail(), bindingResult.getAllErrors());
            return "home/login";
        }

        try {
            UsernamePasswordAuthenticationToken authToken =
                    new UsernamePasswordAuthenticationToken(loginDTO.getEmail(), loginDTO.getPassword());
            Authentication auth = authenticationManager.authenticate(authToken);
            SecurityContextHolder.getContext().setAuthentication(auth);

            log.info("User successfully authenticated: {}", loginDTO.getEmail());
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("home.login.success"),
                    "type", WebUtils.MSG_SUCCESS
            ));
            return "redirect:/dashboard";
        } catch (Exception ex) {
            log.error("Authentication failed for email {}: {}", loginDTO.getEmail(), ex.getMessage());
            redirectAttributes.addFlashAttribute("info", Map.of(
                    "message", WebUtils.getMessage("home.login.failed"),
                    "type", WebUtils.MSG_ERROR
            ));
            return "redirect:/login";
        }
    }
}

