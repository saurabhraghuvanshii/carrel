#include <queue>
#include <string>
#include <vector>
using namespace std;

// Simulate in rounds of n + 1 slots, always taking the most remaining tasks.
int fewestUnits(string& tasks, int n) {
    vector<int> count(26);
    for (char c : tasks) count[c - 'A']++;
    priority_queue<int> heap;
    for (int c : count)
        if (c > 0) heap.push(c);
    int time = 0;
    while (!heap.empty()) {
        vector<int> used;
        int slots = n + 1;
        while (slots > 0 && !heap.empty()) {
            used.push_back(heap.top() - 1);
            heap.pop();
            slots--;
        }
        for (int c : used)
            if (c > 0) heap.push(c);
        time += heap.empty() ? (int)used.size() : n + 1;
    }
    return time;
}
