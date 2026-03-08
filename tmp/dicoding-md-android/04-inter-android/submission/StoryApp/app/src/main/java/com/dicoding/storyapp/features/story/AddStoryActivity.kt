package com.dicoding.storyapp.features.story

import android.Manifest
import android.content.Intent
import android.content.Intent.ACTION_GET_CONTENT
import android.content.pm.PackageManager
import android.graphics.BitmapFactory
import android.net.Uri
import android.os.Bundle
import android.os.Handler
import android.provider.MediaStore
import android.view.View
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.core.content.FileProvider
import com.dicoding.storyapp.common.ResultState
import com.dicoding.storyapp.common.createTempFile
import com.dicoding.storyapp.common.reduceFileImage
import com.dicoding.storyapp.common.showDialog
import com.dicoding.storyapp.common.uriToFile
import com.dicoding.storyapp.databinding.ActivityAddStoryBinding
import com.dicoding.storyapp.features.ViewModelFactory
import com.dicoding.storyapp.features.home.HomeActivity
import java.io.File

class AddStoryActivity : AppCompatActivity() {

    private lateinit var binding: ActivityAddStoryBinding
    private lateinit var currentPhotoPath: String

    private val viewModel by viewModels<StoryViewModel> {
        ViewModelFactory.getInstance(applicationContext)
    }

    private var getFile: File? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityAddStoryBinding.inflate(layoutInflater)
        setContentView(binding.root)

        if (!allPermissionsGranted()) {
            ActivityCompat.requestPermissions(this, REQUIRED_PERMISSIONS, REQUEST_CODE_PERMISSIONS)
        }

        binding.cameraButton.setOnClickListener { startTakePhoto() }
        binding.galleryButton.setOnClickListener { startGallery() }
        binding.uploadButton.setOnClickListener { upload() }
    }

    private fun upload() {
        if (getFile != null) {
            val file = reduceFileImage(getFile as File)
            val description = binding.etDescription.text.toString()
            if (description == "") {
                binding.etDescription.error = "Deskripsi tidak boleh kosong!"
            }

            viewModel.addStory(
                description = description,
                photo = file,
                lat = 0.00f,
                lon = 0.00f
            ).observe(this) {
                when (it) {
                    is ResultState.Loading -> {
                        showLoading(true)
                    }

                    is ResultState.Error -> {
                        showLoading(false)
                        showDialog(
                            this,
                            "Peringatan",
                            "Story tidak bisa diupload."
                        ) { _, _ ->
                            finish()
                        }
                    }

                    is ResultState.Success -> {
                        showLoading(false)
                        val moveIntent = Intent(this, HomeActivity::class.java)
                        startActivity(moveIntent)
                        finish()
                    }
                }
            }

        }
    }

    private fun startTakePhoto() {
        val intent = Intent(MediaStore.ACTION_IMAGE_CAPTURE)
        intent.resolveActivity(packageManager)
        createTempFile(application).also {
            val photoURI: Uri = FileProvider.getUriForFile(
                this,
                "com.dicoding.storyapp",
                it
            )
            currentPhotoPath = it.absolutePath
            intent.putExtra(MediaStore.EXTRA_OUTPUT, photoURI)
            launcherIntentCamera.launch(intent)
        }
    }

    private fun startGallery() {
        val intent = Intent()
        intent.action = ACTION_GET_CONTENT
        intent.type = "image/*"
        val chooser = Intent.createChooser(intent, "Pilih gambar")
        launcherIntentGallery.launch(chooser)
    }

    private val launcherIntentCamera = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) {
        if (it.resultCode == RESULT_OK) {
            val myFile = File(currentPhotoPath)
            getFile = myFile
            val result = BitmapFactory.decodeFile(getFile?.path)
            binding.previewImageView.setImageBitmap(result)
        }
    }


    private val launcherIntentGallery = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == RESULT_OK) {
            val selectedImg: Uri = result.data?.data as Uri
            val myFile = uriToFile(selectedImg, this)
            getFile = myFile
            binding.previewImageView.setImageURI(selectedImg)
        }
    }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<String>,
        grantResults: IntArray
    ) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == REQUEST_CODE_PERMISSIONS) {
            if (!allPermissionsGranted()) {
                Toast.makeText(
                    this,
                    "Tidak mendapatkan permission.",
                    Toast.LENGTH_SHORT
                ).show()
                finish()
            }
        }
    }

    private fun allPermissionsGranted() = REQUIRED_PERMISSIONS.all {
        ContextCompat.checkSelfPermission(baseContext, it) == PackageManager.PERMISSION_GRANTED
    }

    override fun onSupportNavigateUp(): Boolean {
        onBackPressed()
        val backTo = Intent(this, HomeActivity::class.java)
        startActivity(backTo)
        return true
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
    }

    companion object {
        private val REQUIRED_PERMISSIONS = arrayOf(Manifest.permission.CAMERA)
        private const val REQUEST_CODE_PERMISSIONS = 10
    }

}