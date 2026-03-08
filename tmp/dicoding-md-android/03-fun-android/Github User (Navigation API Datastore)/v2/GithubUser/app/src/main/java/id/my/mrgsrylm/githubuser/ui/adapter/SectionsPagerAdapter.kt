package id.my.mrgsrylm.githubuser.ui.adapter

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.viewpager2.adapter.FragmentStateAdapter
import id.my.mrgsrylm.githubuser.ui.fragment.FollowersFragment
import id.my.mrgsrylm.githubuser.ui.fragment.FollowingFragment

class SectionsPagerAdapter(
    activity: AppCompatActivity,
    private val username: Bundle
) : FragmentStateAdapter(activity) {
    override fun getItemCount(): Int {
        return 2
    }

    override fun createFragment(position: Int): Fragment {
        var fragment: Fragment? = null
        when (position) {
            0 -> fragment = FollowersFragment()
            1 -> fragment = FollowingFragment()
        }

        fragment?.arguments = username
        return fragment as Fragment
    }
}