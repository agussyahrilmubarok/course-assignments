package id.my.mrgsrylm.githubuser.view.fragment

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.recyclerview.widget.LinearLayoutManager
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.databinding.FragmentFollowingBinding
import id.my.mrgsrylm.githubuser.view.activity.UserDetailActivity
import id.my.mrgsrylm.githubuser.view.adapter.ListFollowingAdapter
import id.my.mrgsrylm.githubuser.viewmodel.FollowingViewModel

class FollowingFragment : Fragment() {
    private val followingViewModel by viewModels<FollowingViewModel>()
    private lateinit var binding: FragmentFollowingBinding

    companion object {
        const val TAG = "FragmentFollowing"
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentFollowingBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        followingViewModel.findFollowing(
            arguments?.getString(UserDetailActivity.EXTRA_FRAGMENT).toString()
        )

        followingViewModel.following.observe(viewLifecycleOwner) { following ->
            setFollowingList(following)
        }

        followingViewModel.isLoading.observe(viewLifecycleOwner) {
            showLoading(it)
        }
    }

    private fun setFollowingList(following: List<UserItemResponse>) {
        val users = ArrayList<UserItemResponse>()
        for (i in following) {
            users.clear()
            users.addAll(following)
        }
        binding.rvFollowing.layoutManager = LinearLayoutManager(requireActivity())
        val adapter = ListFollowingAdapter(users)
        binding.rvFollowing.adapter = adapter
    }

    private fun showLoading(isLoading: Boolean) {
        if (isLoading) {
            binding.progressBar.visibility = View.VISIBLE
        } else {
            binding.progressBar.visibility = View.GONE
        }
    }
}