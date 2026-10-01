#include <deque>
#include <unordered_map>
#include <utility>
using namespace std;

class VisitorCounter {
    int window;
    deque<pair<int, int>> visits;  // {user, time}
    unordered_map<int, int> latest;

public:
    VisitorCounter(int window) : window(window) {}

    void visit(int user, int time) {
        visits.push_back({user, time});
        latest[user] = time;
    }

    int unique(int time) {
        while (!visits.empty() && visits.front().second <= time - window) {
            auto [user, at] = visits.front();
            visits.pop_front();
            // The user may be gone already if they visited twice at this time.
            auto it = latest.find(user);
            if (it != latest.end() && it->second == at) latest.erase(it);
        }
        return latest.size();
    }
};
