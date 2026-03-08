package com.dicoding.asclepius.view

import android.net.Uri
import android.os.Bundle
import android.util.Log
import androidx.appcompat.app.AppCompatActivity
import com.dicoding.asclepius.R
import com.dicoding.asclepius.databinding.ActivityResultBinding

class ResultActivity : AppCompatActivity() {
    private lateinit var binding: ActivityResultBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_result)
        binding = ActivityResultBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val imageUri = Uri.parse(intent.getStringExtra(EXTRA_IMAGE_URI))
        imageUri?.let {
            Log.d(EXTRA_IMAGE_URI, "showImage: $it")
            binding.resultImage.setImageURI(it)
        }

        val result = intent.getStringExtra(EXTRA_RESULT)
        result?.let {
            Log.d(EXTRA_RESULT, "showResult: $it")
            binding.resultText.text = it
        }
    }

    companion object {
        const val EXTRA_IMAGE_URI = "image_uri"
        const val EXTRA_RESULT = "result"
    }
}