#include <algorithm>
#include <deque>
#include <unordered_map>
#include <utility>
#include <vector>
using namespace std;

class SeatHolder {
    struct Hold {
        int first, count, until;
    };
    vector<bool> taken;  // held or sold
    int ttl, freeCount;
    unordered_map<int, Hold> holds;
    deque<pair<int, int>> byExpiry;  // {user, until}

    void giveBack(int user) {
        Hold h = holds[user];
        holds.erase(user);
        fill(taken.begin() + h.first, taken.begin() + h.first + h.count, false);
        freeCount += h.count;
    }

    void expire(int time) {
        while (!byExpiry.empty() && byExpiry.front().second <= time) {
            auto [user, until] = byExpiry.front();
            byExpiry.pop_front();
            auto it = holds.find(user);
            if (it != holds.end() && it->second.until == until) giveBack(user);
        }
    }

public:
    SeatHolder(int seats, int ttl) : taken(seats, false), ttl(ttl), freeCount(seats) {}

    int hold(int user, int count, int time) {
        expire(time);
        if (holds.count(user)) return -1;
        int run = 0;
        for (int s = 0; s < (int)taken.size(); s++) {
            run = taken[s] ? 0 : run + 1;
            if (run == count) {
                int first = s - count + 1;
                fill(taken.begin() + first, taken.begin() + s + 1, true);
                freeCount -= count;
                holds[user] = {first, count, time + ttl};
                byExpiry.push_back({user, time + ttl});
                return first;
            }
        }
        return -1;
    }

    bool buy(int user, int time) {
        expire(time);
        return holds.erase(user) > 0;
    }

    bool release(int user, int time) {
        expire(time);
        if (!holds.count(user)) return false;
        giveBack(user);
        return true;
    }

    int freeSeats(int time) {
        expire(time);
        return freeCount;
    }
};
