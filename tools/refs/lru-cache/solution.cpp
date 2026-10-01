#include <list>
#include <unordered_map>
#include <utility>
using namespace std;

class LRUCache {
    int capacity;
    list<pair<int, int>> order;
    unordered_map<int, list<pair<int, int>>::iterator> byKey;

public:
    LRUCache(int capacity) : capacity(capacity) {}

    int get(int key) {
        auto it = byKey.find(key);
        if (it == byKey.end()) return -1;
        order.splice(order.begin(), order, it->second);
        return it->second->second;
    }

    void put(int key, int value) {
        auto it = byKey.find(key);
        if (it != byKey.end()) {
            it->second->second = value;
            order.splice(order.begin(), order, it->second);
            return;
        }
        order.emplace_front(key, value);
        byKey[key] = order.begin();
        if ((int)order.size() > capacity) {
            byKey.erase(order.back().first);
            order.pop_back();
        }
    }
};
