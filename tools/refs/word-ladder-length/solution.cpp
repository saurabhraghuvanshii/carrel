#include <queue>
#include <string>
#include <unordered_set>
#include <vector>
using namespace std;

int ladderLength(string& start, string& end, vector<string>& words) {
    unordered_set<string> left(words.begin(), words.end());
    if (!left.count(end)) return 0;
    queue<string> q;
    q.push(start);
    left.erase(start);
    for (int steps = 1; !q.empty(); steps++) {
        for (int size = q.size(); size > 0; size--) {
            string w = q.front();
            q.pop();
            if (w == end) return steps;
            for (size_t i = 0; i < w.size(); i++) {
                char keep = w[i];
                for (char ch = 'a'; ch <= 'z'; ch++) {
                    w[i] = ch;
                    if (left.erase(w)) q.push(w);
                }
                w[i] = keep;
            }
        }
    }
    return 0;
}
