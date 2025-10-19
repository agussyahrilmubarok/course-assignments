package com.dicoding.storyapp.features.register

import android.animation.AnimatorSet
import android.animation.ObjectAnimator
import android.content.Intent
import android.os.Build
import android.os.Bundle
import android.os.Handler
import android.view.View
import android.view.WindowInsets
import android.view.WindowManager
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import com.dicoding.storyapp.common.InputValidator
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.common.showDialog
import com.dicoding.storyapp.data.model.RegisterRequest
import com.dicoding.storyapp.databinding.ActivityRegisterBinding
import com.dicoding.storyapp.features.ViewModelFactory
import com.dicoding.storyapp.features.login.LoginActivity

class RegisterActivity : AppCompatActivity() {

    private lateinit var binding: ActivityRegisterBinding

    private val viewModel by viewModels<RegisterViewModel> {
        ViewModelFactory.getInstance(applicationContext)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityRegisterBinding.inflate(layoutInflater)
        setContentView(binding.root)
        setupView()
        setupAction()
        playAnimation()
    }

    private fun setupView() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
            window.insetsController?.hide(WindowInsets.Type.statusBars())
        } else {
            window.setFlags(
                WindowManager.LayoutParams.FLAG_FULLSCREEN,
                WindowManager.LayoutParams.FLAG_FULLSCREEN
            )
        }
        supportActionBar?.hide()
    }

    private fun setupAction() {
        binding.signupButton.isEnabled = true
        binding.signupButton.setOnClickListener {
            val name = binding.edRegisterName.text.toString()
            val email = binding.edRegisterEmail.text.toString()
            val password = binding.edRegisterPassword.text.toString()

            if (validateInput(name, email)) {
                doRegister(RegisterRequest(name, email, password))
            }
        }
    }

    private fun playAnimation() {
        ObjectAnimator.ofFloat(binding.imageView, View.TRANSLATION_X, -30f, 30f).apply {
            duration = 6000
            repeatCount = ObjectAnimator.INFINITE
            repeatMode = ObjectAnimator.REVERSE
        }.start()

        val title = ObjectAnimator.ofFloat(binding.titleTextView, View.ALPHA, 1f).setDuration(100)
        val nameTextView =
            ObjectAnimator.ofFloat(binding.nameTextView, View.ALPHA, 1f).setDuration(100)
        val nameEditTextLayout =
            ObjectAnimator.ofFloat(binding.nameEditTextLayout, View.ALPHA, 1f).setDuration(100)
        val emailTextView =
            ObjectAnimator.ofFloat(binding.emailTextView, View.ALPHA, 1f).setDuration(100)
        val emailEditTextLayout =
            ObjectAnimator.ofFloat(binding.emailEditTextLayout, View.ALPHA, 1f).setDuration(100)
        val passwordTextView =
            ObjectAnimator.ofFloat(binding.passwordTextView, View.ALPHA, 1f).setDuration(100)
        val passwordEditTextLayout =
            ObjectAnimator.ofFloat(binding.passwordEditTextLayout, View.ALPHA, 1f).setDuration(100)
        val signup = ObjectAnimator.ofFloat(binding.signupButton, View.ALPHA, 1f).setDuration(100)


        AnimatorSet().apply {
            playSequentially(
                title,
                nameTextView,
                nameEditTextLayout,
                emailTextView,
                emailEditTextLayout,
                passwordTextView,
                passwordEditTextLayout,
                signup
            )
            startDelay = 100
        }.start()
    }

    private fun doRegister(param: RegisterRequest) {
        viewModel.register(param).observe(this) {
            when (it) {
                is ResultState.Loading -> {
                    binding.signupButton.isEnabled = false
                }

                is ResultState.Error -> {
                    binding.signupButton.isEnabled = true
                    showDialog(this, "Register Gagal", it.message)
                }

                is ResultState.Success -> {
                    binding.signupButton.isEnabled = true
                    showDialog(this, "Register Berhasil", it.data.message)
                    Handler().postDelayed({
                        val moveIntent = Intent(this, LoginActivity::class.java)
                        startActivity(moveIntent)
                        finish()
                    }, 3000)
                }
            }
        }
    }

    private fun validateInput(name: String, email: String): Boolean {
        var isValid = true

        if (!InputValidator.isValidName(name)) {
            binding.edRegisterName.error = "Name is required"
            isValid = false
        } else {
            binding.edRegisterName.error = null
        }

        if (!InputValidator.isValidEmail(email)) {
            binding.edRegisterEmail.error = "Invalid email address"
            isValid = false
        } else {
            binding.edRegisterEmail.error = null
        }

        return isValid
    }
}