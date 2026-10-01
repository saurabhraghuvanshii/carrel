#include <unordered_map>
#include <utility>
using namespace std;

class RateLimiter {
    int limit, window;
    unordered_map<int, pair<int, int>> seen;

public:
    RateLimiter(int limit, int window) : limit(limit), window(window) {}

    bool allow(int user, int time) {
        int w = time / window;
        auto it = seen.find(user);
        if (it == seen.end() || it->second.first != w) {
            seen[user] = {w, 1};
            return true;
        }
        if (it->second.second == limit) return false;
        it->second.second++;
        return true;
    }
};
