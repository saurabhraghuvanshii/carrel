#include <string>
using namespace std;

string removeTwins(string& text) {
    string kept;
    for (char ch : text) {
        if (!kept.empty() && kept.back() == ch) kept.pop_back();
        else kept.push_back(ch);
    }
    return kept;
}
