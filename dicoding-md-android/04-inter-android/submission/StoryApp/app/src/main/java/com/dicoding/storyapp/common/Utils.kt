package com.dicoding.storyapp.common

import android.app.Activity
import android.app.Application
import android.content.ContentResolver
import android.content.Context
import android.content.DialogInterface
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.net.Uri
import android.os.Environment
import androidx.appcompat.app.AlertDialog
import com.dicoding.storyapp.R
import java.io.ByteArrayOutputStream
import java.io.File
import java.io.FileOutputStream
import java.io.InputStream
import java.io.OutputStream
import java.text.SimpleDateFormat
import java.util.Locale

fun showDialog(ctx: Context, title: String, message: String) {
    AlertDialog.Builder(ctx).apply {
        setTitle(title)
        setMessage(message)
        setPositiveButton("OK") { _, _ ->
            if (ctx is Activity) {
                ctx.finish()
            }
        }
        setCancelable(false)
        create()
        show()
    }
}

fun showDialog(
    ctx: Context,
    title: String,
    message: String,
    onOk: (dialog: DialogInterface, which: Int) -> Unit
) {
    AlertDialog.Builder(ctx).apply {
        setTitle(title)
        setMessage(message)
        setCancelable(false)
        setPositiveButton("OK", onOk)
        setCancelable(false)
        create()
        show()
    }
}

fun showNoInternet(ctx: Context, onRetry: (dialog: DialogInterface, which: Int) -> Unit) {
    AlertDialog.Builder(ctx).apply {
        setTitle("Peringatan")
        setMessage("Anda tidak memiliki koneksi internet")
        setPositiveButton("OK", onRetry)
        setCancelable(false)
        create()
        show()
    }
}

private const val FILENAME_FORMAT = "dd-MMM-yyyy"

val timeStamp: String = SimpleDateFormat(
    FILENAME_FORMAT,
    Locale.US
).format(System.currentTimeMillis())

fun createTempFile(context: Context): File {
    val storageDir: File? = context.getExternalFilesDir(Environment.DIRECTORY_PICTURES)
    return File.createTempFile(timeStamp, ".jpg", storageDir)
}

fun createFile(application: Application): File {
    val mediaDir = application.externalMediaDirs.firstOrNull()?.let {
        File(it, application.resources.getString(R.string.app_name)).apply { mkdirs() }
    }

    val outputDirectory = if (
        mediaDir != null && mediaDir.exists()
    ) mediaDir else application.filesDir

    return File(outputDirectory, "$timeStamp.jpg")
}

fun uriToFile(selectedImg: Uri, context: Context): File {
    val contentResolver: ContentResolver = context.contentResolver
    val myFile = createTempFile(context)

    val inputStream = contentResolver.openInputStream(selectedImg) as InputStream
    val outputStream: OutputStream = FileOutputStream(myFile)
    val buf = ByteArray(1024)
    var len: Int
    while (inputStream.read(buf).also { len = it } > 0) outputStream.write(buf, 0, len)
    outputStream.close()
    inputStream.close()

    return myFile
}

fun reduceFileImage(file: File): File {
    val bitmap = BitmapFactory.decodeFile(file.path)
    var compressQuality = 100
    var streamLength: Int
    do {
        val bmpStream = ByteArrayOutputStream()
        bitmap.compress(Bitmap.CompressFormat.JPEG, compressQuality, bmpStream)
        val bmpPicByteArray = bmpStream.toByteArray()
        streamLength = bmpPicByteArray.size
        compressQuality -= 5
    } while (streamLength > 1000000)
    bitmap.compress(Bitmap.CompressFormat.JPEG, compressQuality, FileOutputStream(file))
    return file
}
