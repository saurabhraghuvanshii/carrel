#include <functional>
#include <queue>
#include <unordered_map>
#include <utility>
#include <vector>
using namespace std;

class TTLCache {
    unordered_map<int, pair<int, int>> live;  // key -> {value, expiry}
    priority_queue<pair<int, int>, vector<pair<int, int>>, greater<>> byExpiry;

public:
    TTLCache() {}

    void put(int key, int value, int ttl, int time) {
        live[key] = {value, time + ttl};
        byExpiry.push({time + ttl, key});
    }

    int get(int key, int time) {
        auto it = live.find(key);
        return it != live.end() && it->second.second > time ? it->second.first : -1;
    }

    int count(int time) {
        while (!byExpiry.empty() && byExpiry.top().first <= time) {
            auto [at, key] = byExpiry.top();
            byExpiry.pop();
            auto it = live.find(key);
            if (it != live.end() && it->second.second == at) live.erase(it);
        }
        return live.size();
    }
};
