#include <deque>
#include <unordered_map>
using namespace std;

class SlidingLimiter {
    int limit, window;
    unordered_map<int, deque<int>> recent;

public:
    SlidingLimiter(int limit, int window) : limit(limit), window(window) {}

    bool allow(int user, int time) {
        auto& q = recent[user];
        while (!q.empty() && q.front() <= time - window) q.pop_front();
        if ((int)q.size() == limit) return false;
        q.push_back(time);
        return true;
    }
};
