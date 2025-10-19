package id.my.mrgsrylm.githubuser.di

import androidx.test.espresso.IdlingResource
import androidx.test.espresso.idling.CountingIdlingResource

object Idling {

    private val RES = "GLOBAL"
    private val countingIdlingResource = CountingIdlingResource(RES)

    val idlingResource: IdlingResource
        get() = countingIdlingResource

    fun increment() {
        countingIdlingResource.increment()
    }

    fun decrement() {
        countingIdlingResource.decrement()
    }
}