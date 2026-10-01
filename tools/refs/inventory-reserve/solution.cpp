#include <unordered_map>
#include <utility>
using namespace std;

class Inventory {
    unordered_map<int, int> freeUnits;
    unordered_map<int, pair<int, int>> held;  // order -> {item, qty}

public:
    Inventory() {}

    void add(int item, int qty) { freeUnits[item] += qty; }

    bool reserve(int order, int item, int qty) {
        if (held.count(order) || freeUnits[item] < qty) return false;
        freeUnits[item] -= qty;
        held[order] = {item, qty};
        return true;
    }

    bool cancel(int order) {
        auto it = held.find(order);
        if (it == held.end()) return false;
        freeUnits[it->second.first] += it->second.second;
        held.erase(it);
        return true;
    }

    bool confirm(int order) { return held.erase(order) > 0; }

    int available(int item) { return freeUnits[item]; }
};
