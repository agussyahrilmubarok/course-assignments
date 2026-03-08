package id.my.mrgsrylm.footballlovers.features.team

import android.content.Context
import android.content.Intent
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.bumptech.glide.request.RequestOptions
import id.my.mrgsrylm.footballlovers.R
import id.my.mrgsrylm.footballlovers.data.model.Team

class TeamAdapter(
    private val items: List<Team>
) : RecyclerView.Adapter<TeamAdapter.ListViewHolder>() {

    override fun onCreateViewHolder(
        parent: ViewGroup,
        viewType: Int
    ): ListViewHolder {
        val view: View = LayoutInflater
            .from(parent.context)
            .inflate(R.layout.layout_team_tile, parent, false)

        return ListViewHolder(view)
    }

    override fun getItemCount(): Int {
        return items.size
    }

    override fun onBindViewHolder(holder: ListViewHolder, position: Int) {
        val team = items[position]

        Glide.with(holder.itemView.context)
            .load(team.imageUrl)
            .apply(RequestOptions())
            .into(holder.ivImageUrl)

        holder.tvName.text = team.name
        holder.tvLeague.text = team.league

        holder.itemView.setOnClickListener {
            val mContext: Context = holder.itemView.context
            val moveDetail = Intent(mContext, TeamActivity::class.java).apply {
                putExtra(TeamActivity.EXTRA_NAME, team.name)
                putExtra(TeamActivity.EXTRA_IMAGE_URL, team.imageUrl)
                putExtra(TeamActivity.EXTRA_DESCRIPTION, team.description)
                putExtra(TeamActivity.EXTRA_STADIUM, team.stadium)
                putExtra(TeamActivity.EXTRA_FOUNDED_YEAR, team.foundedYear)
                putExtra(TeamActivity.EXTRA_COACH, team.coach)
                putExtra(TeamActivity.EXTRA_COUNTRY, team.country)
                putExtra(TeamActivity.EXTRA_LEAGUE, team.league)
            }
            mContext.startActivity(moveDetail)
        }
    }

    inner class ListViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        var tvName: TextView = itemView.findViewById(R.id.tv_set_team_name)
        var tvLeague: TextView = itemView.findViewById(R.id.tv_set_team_league)
        var ivImageUrl: ImageView = itemView.findViewById(R.id.iv_set_team_image_url)
    }
}