#include <list>
#include <unordered_map>
using namespace std;

class LFUCache {
    struct Entry {
        int value, uses;
        list<int>::iterator pos;
    };
    int capacity, lowest = 0;
    unordered_map<int, Entry> entries;
    unordered_map<int, list<int>> byUses;  // uses -> keys, oldest use first

    void use(int key, Entry& e) {
        auto& l = byUses[e.uses];
        l.erase(e.pos);
        if (l.empty() && lowest == e.uses) lowest++;
        e.uses++;
        auto& next = byUses[e.uses];
        e.pos = next.insert(next.end(), key);
    }

public:
    LFUCache(int capacity) : capacity(capacity) {}

    int get(int key) {
        auto it = entries.find(key);
        if (it == entries.end()) return -1;
        use(key, it->second);
        return it->second.value;
    }

    void put(int key, int value) {
        auto it = entries.find(key);
        if (it != entries.end()) {
            it->second.value = value;
            use(key, it->second);
            return;
        }
        if ((int)entries.size() == capacity) {
            auto& l = byUses[lowest];
            entries.erase(l.front());
            l.pop_front();
        }
        auto& ones = byUses[1];
        entries[key] = {value, 1, ones.insert(ones.end(), key)};
        lowest = 1;
    }
};
